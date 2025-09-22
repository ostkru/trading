-- Интеграция существующей таблицы tariffs с системой rate limiting
-- Файл: add_tariff_integration.sql

USE portaldata;

-- 1. Добавляем поле tariff_id в таблицу users
ALTER TABLE users ADD COLUMN tariff_id INT DEFAULT 1;

-- 2. Добавляем внешний ключ на таблицу tariffs
ALTER TABLE users ADD FOREIGN KEY (tariff_id) REFERENCES tariffs(id) ON DELETE SET NULL;

-- 3. Создаем индекс для оптимизации
CREATE INDEX idx_users_tariff_id ON users(tariff_id);

-- 4. Обновляем существующих пользователей - назначаем им базовый тариф
UPDATE users SET tariff_id = 1 WHERE tariff_id IS NULL;

-- 5. Создаем функцию для получения лимитов пользователя
DELIMITER //
CREATE FUNCTION GetUserDailyLimit(user_id INT) 
RETURNS INT
READS SQL DATA
DETERMINISTIC
BEGIN
    DECLARE daily_limit INT DEFAULT 1000;
    
    SELECT JSON_UNQUOTE(JSON_EXTRACT(features, '$.daily_requests_limit'))
    INTO daily_limit
    FROM users u
    JOIN tariffs t ON u.tariff_id = t.id
    WHERE u.id = user_id AND u.is_active = 1 AND t.is_active = 1;
    
    -- Если лимит не найден или NULL, возвращаем дефолтный
    IF daily_limit IS NULL OR daily_limit <= 0 THEN
        SET daily_limit = 1000;
    END IF;
    
    RETURN daily_limit;
END //
DELIMITER ;

-- 6. Создаем функцию для получения минутного лимита пользователя
DELIMITER //
CREATE FUNCTION GetUserMinuteLimit(user_id INT) 
RETURNS INT
READS SQL DATA
DETERMINISTIC
BEGIN
    DECLARE minute_limit INT DEFAULT 60;
    DECLARE daily_limit INT;
    
    -- Получаем дневной лимит
    SELECT GetUserDailyLimit(user_id) INTO daily_limit;
    
    -- Минутный лимит = дневной лимит / 24 / 60 (примерно)
    -- Но минимум 60 запросов в минуту
    SET minute_limit = GREATEST(60, daily_limit / 1440);
    
    RETURN minute_limit;
END //
DELIMITER ;

-- 7. Создаем представление для удобного просмотра пользователей с тарифами и лимитами
CREATE OR REPLACE VIEW user_tariffs_with_limits AS
SELECT 
    u.id as user_id,
    u.username,
    u.email,
    u.api_token,
    u.is_active as user_active,
    u.created_at as user_created_at,
    t.id as tariff_id,
    t.name as tariff_name,
    t.price as tariff_price,
    t.duration_days,
    t.is_active as tariff_active,
    JSON_UNQUOTE(JSON_EXTRACT(t.features, '$.daily_requests_limit')) as daily_limit,
    GetUserMinuteLimit(u.id) as minute_limit,
    GetUserDailyLimit(u.id) as calculated_daily_limit
FROM users u
LEFT JOIN tariffs t ON u.tariff_id = t.id;

-- 8. Создаем процедуру для безопасного изменения тарифа пользователя
DELIMITER //
CREATE PROCEDURE ChangeUserTariff(
    IN p_user_id INT,
    IN p_tariff_id INT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- Проверяем существование пользователя
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = p_user_id AND is_active = 1) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Пользователь не найден или неактивен';
    END IF;
    
    -- Проверяем существование тарифа
    IF NOT EXISTS (SELECT 1 FROM tariffs WHERE id = p_tariff_id AND is_active = 1) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Тариф не найден или неактивен';
    END IF;
    
    -- Обновляем тариф пользователя
    UPDATE users 
    SET tariff_id = p_tariff_id, updated_at = NOW() 
    WHERE id = p_user_id;
    
    -- Сбрасываем счетчики rate limiting для пользователя
    UPDATE api_rate_limits 
    SET minute_count = 0, day_count = 0, updated_at = NOW() 
    WHERE user_id = p_user_id;
    
    COMMIT;
END //
DELIMITER ;

-- 9. Создаем процедуру для получения лимитов пользователя
DELIMITER //
CREATE PROCEDURE GetUserLimits(
    IN p_user_id INT,
    OUT p_minute_limit INT,
    OUT p_daily_limit INT
)
BEGIN
    SET p_daily_limit = GetUserDailyLimit(p_user_id);
    SET p_minute_limit = GetUserMinuteLimit(p_user_id);
END //
DELIMITER ;

-- 10. Создаем представление для статистики использования тарифов
CREATE OR REPLACE VIEW tariff_usage_statistics AS
SELECT 
    t.id as tariff_id,
    t.name as tariff_name,
    t.price as tariff_price,
    JSON_UNQUOTE(JSON_EXTRACT(t.features, '$.daily_requests_limit')) as daily_limit,
    COUNT(u.id) as user_count,
    COALESCE(AVG(arl.day_count), 0) as avg_daily_usage,
    COALESCE(MAX(arl.last_request_time), NULL) as last_activity,
    t.is_active as tariff_active
FROM tariffs t
LEFT JOIN users u ON t.id = u.tariff_id AND u.is_active = 1
LEFT JOIN api_rate_limits arl ON u.id = arl.user_id
WHERE t.is_active = 1
GROUP BY t.id, t.name, t.price, t.features, t.is_active
ORDER BY t.price ASC;

-- 11. Добавляем тестовых пользователей с разными тарифами
INSERT INTO users (username, name, email, api_token, tariff_id, is_active) VALUES
('test_user_1', 'Тестовый пользователь 1', 'test1@example.com', 'test_api_key_1', 1, 1),
('test_user_2', 'Тестовый пользователь 2', 'test2@example.com', 'test_api_key_2', 2, 1),
('test_user_3', 'Тестовый пользователь 3', 'test3@example.com', 'test_api_key_3', 3, 1),
('test_user_4', 'Тестовый пользователь 4', 'test4@example.com', 'test_api_key_4', 4, 1),
('test_user_5', 'Тестовый пользователь 5', 'test5@example.com', 'test_api_key_5', 5, 1)
ON DUPLICATE KEY UPDATE 
    tariff_id = VALUES(tariff_id),
    is_active = VALUES(is_active);

-- 12. Выводим информацию о созданных функциях и процедурах
SELECT 'Интеграция тарифов завершена' as status;
SELECT 'Созданы функции:' as info;
SELECT '  - GetUserDailyLimit(user_id)' as function;
SELECT '  - GetUserMinuteLimit(user_id)' as function;
SELECT 'Созданы процедуры:' as info;
SELECT '  - ChangeUserTariff(user_id, tariff_id)' as procedure;
SELECT '  - GetUserLimits(user_id, @minute_limit, @daily_limit)' as procedure;
SELECT 'Созданы представления:' as info;
SELECT '  - user_tariffs_with_limits' as view;
SELECT '  - tariff_usage_statistics' as view;

-- 13. Показываем текущие тарифы с лимитами
SELECT 
    id,
    name,
    price,
    JSON_UNQUOTE(JSON_EXTRACT(features, '$.daily_requests_limit')) as daily_limit,
    is_active
FROM tariffs 
ORDER BY price ASC;
