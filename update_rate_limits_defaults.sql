-- Обновление таблицы api_rate_limits с данными по умолчанию
-- Добавляем записи для существующих пользователей с лимитами по умолчанию

-- Очищаем таблицу (если есть данные)
TRUNCATE TABLE api_rate_limits;

-- Добавляем записи для всех пользователей с лимитами по умолчанию
-- 60 запросов в минуту и 1000 в сутки для всех методов
INSERT INTO api_rate_limits (user_id, api_key, endpoint, request_count, minute_count, day_count, last_request_time, minute_start, day_start, created_at, updated_at)
SELECT 
    u.id as user_id,
    u.api_token,
    'all' as endpoint,
    0 as request_count,
    0 as minute_count,
    0 as day_count,
    NOW() as last_request_time,
    NOW() as minute_start,
    NOW() as day_start,
    NOW() as created_at,
    NOW() as updated_at
FROM users u;

-- Добавляем записи для публичных методов (offers/public)
INSERT INTO api_rate_limits (user_id, api_key, endpoint, request_count, minute_count, day_count, last_request_time, minute_start, day_start, created_at, updated_at)
SELECT 
    u.id as user_id,
    u.api_token,
    'public' as endpoint,
    0 as request_count,
    0 as minute_count,
    0 as day_count,
    NOW() as last_request_time,
    NOW() as minute_start,
    NOW() as day_start,
    NOW() as created_at,
    NOW() as updated_at
FROM users u;

-- Проверяем результат
SELECT 
    user_id,
    api_key,
    endpoint,
    minute_count,
    day_count,
    last_request_time
FROM api_rate_limits
ORDER BY user_id, endpoint; 