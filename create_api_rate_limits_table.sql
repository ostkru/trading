-- Создание таблицы api_rate_limits для проверки лимитов запросов
-- Таблица будет использоваться для отслеживания количества запросов в минуту и день

CREATE TABLE api_rate_limits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    request_count INTEGER DEFAULT 1,
    minute_count INTEGER DEFAULT 1,
    day_count INTEGER DEFAULT 1,
    last_request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    minute_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    day_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Индексы для быстрого поиска
    INDEX idx_user_api_endpoint (user_id, api_key, endpoint),
    INDEX idx_last_request_time (last_request_time),
    INDEX idx_minute_start (minute_start),
    INDEX idx_day_start (day_start)
);

-- Добавляем комментарии к таблице
ALTER TABLE api_rate_limits COMMENT = 'Таблица для отслеживания лимитов API запросов'; 