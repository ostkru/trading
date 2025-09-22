-- Интеграция существующей таблицы tariffs с системой rate limiting (упрощенная версия)
-- Файл: add_tariff_integration_simple.sql

USE portaldata;

-- 1. Добавляем поле tariff_id в таблицу users
ALTER TABLE users ADD COLUMN tariff_id INT DEFAULT 1;

-- 2. Создаем индекс для оптимизации
CREATE INDEX idx_users_tariff_id ON users(tariff_id);

-- 3. Обновляем существующих пользователей - назначаем им базовый тариф
UPDATE users SET tariff_id = 1 WHERE tariff_id IS NULL;

-- 4. Создаем представление для удобного просмотра пользователей с тарифами и лимитами
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
    GREATEST(60, JSON_UNQUOTE(JSON_EXTRACT(t.features, '$.daily_requests_limit')) / 1440) as minute_limit
FROM users u
LEFT JOIN tariffs t ON u.tariff_id = t.id;

-- 5. Создаем представление для статистики использования тарифов
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

-- 6. Добавляем тестовых пользователей с разными тарифами
INSERT INTO users (username, name, email, api_token, tariff_id, is_active) VALUES
('test_user_1', 'Тестовый пользователь 1', 'test1@example.com', 'test_api_key_1', 1, 1),
('test_user_2', 'Тестовый пользователь 2', 'test2@example.com', 'test_api_key_2', 2, 1),
('test_user_3', 'Тестовый пользователь 3', 'test3@example.com', 'test_api_key_3', 3, 1),
('test_user_4', 'Тестовый пользователь 4', 'test4@example.com', 'test_api_key_4', 4, 1),
('test_user_5', 'Тестовый пользователь 5', 'test5@example.com', 'test_api_key_5', 5, 1)
ON DUPLICATE KEY UPDATE 
    tariff_id = VALUES(tariff_id),
    is_active = VALUES(is_active);

-- 7. Выводим информацию о созданных представлениях
SELECT 'Интеграция тарифов завершена' as status;
SELECT 'Созданы представления:' as info;
SELECT '  - user_tariffs_with_limits' as view;
SELECT '  - tariff_usage_statistics' as view;

-- 8. Показываем текущие тарифы с лимитами
SELECT 
    id,
    name,
    price,
    JSON_UNQUOTE(JSON_EXTRACT(features, '$.daily_requests_limit')) as daily_limit,
    is_active
FROM tariffs 
ORDER BY price ASC;
