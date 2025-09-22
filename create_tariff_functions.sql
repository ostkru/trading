-- Создание функций для работы с тарифами
-- Файл: create_tariff_functions.sql

USE portaldata;

-- Функция для получения дневного лимита пользователя
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
END;

-- Функция для получения минутного лимита пользователя
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
END;

-- Процедура для безопасного изменения тарифа пользователя
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
END;

-- Процедура для получения лимитов пользователя
CREATE PROCEDURE GetUserLimits(
    IN p_user_id INT,
    OUT p_minute_limit INT,
    OUT p_daily_limit INT
)
BEGIN
    SET p_daily_limit = GetUserDailyLimit(p_user_id);
    SET p_minute_limit = GetUserMinuteLimit(p_user_id);
END;
