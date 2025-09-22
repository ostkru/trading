-- Добавление таблицы tariffs и интеграция с пользователями
-- Файл: add_tariffs_table.sql

USE portaldata;

-- Создание таблицы tariffs
CREATE TABLE IF NOT EXISTS tariffs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    minute_limit INT NOT NULL DEFAULT 1000,
    day_limit INT NOT NULL DEFAULT 10000,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Индексы для оптимизации
    INDEX idx_tariffs_active (is_active),
    INDEX idx_tariffs_name (name)
);

-- Создание таблицы users если её нет
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    api_token VARCHAR(255) NOT NULL UNIQUE,
    tariff_id INT DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Индексы для оптимизации
    INDEX idx_users_api_token (api_token),
    INDEX idx_users_tariff_id (tariff_id),
    INDEX idx_users_active (is_active),
    
    -- Внешний ключ на тарифы
    FOREIGN KEY (tariff_id) REFERENCES tariffs(id) ON DELETE SET NULL
);

-- Создание таблицы api_rate_limits если её нет
CREATE TABLE IF NOT EXISTS api_rate_limits (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    endpoint VARCHAR(50) NOT NULL,
    request_count INT DEFAULT 0,
    minute_count INT DEFAULT 0,
    day_count INT DEFAULT 0,
    last_request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    minute_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    day_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Индексы для оптимизации
    INDEX idx_rate_limits_user_api (user_id, api_key),
    INDEX idx_rate_limits_endpoint (endpoint),
    INDEX idx_rate_limits_last_request (last_request_time),
    
    -- Внешний ключ на пользователей
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Вставка базовых тарифов
INSERT INTO tariffs (name, description, minute_limit, day_limit) VALUES
('Базовый', 'Базовый тарифный план для новых пользователей', 1000, 10000),
('Стандартный', 'Стандартный тарифный план с увеличенными лимитами', 3000, 30000),
('Премиум', 'Премиум тарифный план для активных пользователей', 5000, 50000),
('Корпоративный', 'Корпоративный тарифный план с максимальными лимитами', 10000, 100000),
('Безлимитный', 'Безлимитный тарифный план (только для администраторов)', 999999, 9999999)
ON DUPLICATE KEY UPDATE 
    description = VALUES(description),
    minute_limit = VALUES(minute_limit),
    day_limit = VALUES(day_limit);

-- Создание административного пользователя по умолчанию
INSERT INTO users (username, email, api_token, tariff_id) VALUES
('admin', 'admin@portaldata.ru', 'f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53', 5)
ON DUPLICATE KEY UPDATE 
    tariff_id = VALUES(tariff_id);

-- Создание тестовых пользователей
INSERT INTO users (username, email, api_token, tariff_id) VALUES
('test_user_1', 'test1@portaldata.ru', 'test_api_key_1', 1),
('test_user_2', 'test2@portaldata.ru', 'test_api_key_2', 2),
('test_user_3', 'test3@portaldata.ru', 'test_api_key_3', 3)
ON DUPLICATE KEY UPDATE 
    tariff_id = VALUES(tariff_id);

-- Создание индексов для оптимизации производительности
CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id);
CREATE INDEX IF NOT EXISTS idx_offers_user_id ON offers(user_id);
CREATE INDEX IF NOT EXISTS idx_warehouses_user_id ON warehouses(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_initiator_user_id ON orders(initiator_user_id);
CREATE INDEX IF NOT EXISTS idx_orders_counterparty_user_id ON orders(counterparty_user_id);

-- Обновление существующих записей products, если user_id NULL
UPDATE products SET user_id = 1 WHERE user_id IS NULL;

-- Обновление существующих записей offers, если user_id NULL
UPDATE offers SET user_id = 1 WHERE user_id IS NULL;

-- Обновление существующих записей warehouses, если user_id NULL
UPDATE warehouses SET user_id = 1 WHERE user_id IS NULL;

-- Обновление существующих записей orders, если initiator_user_id NULL
UPDATE orders SET initiator_user_id = 1 WHERE initiator_user_id IS NULL;

-- Обновление существующих записей orders, если counterparty_user_id NULL
UPDATE orders SET counterparty_user_id = 1 WHERE counterparty_user_id IS NULL;

-- Создание представления для удобного просмотра пользователей с тарифами
CREATE OR REPLACE VIEW user_tariffs_view AS
SELECT 
    u.id,
    u.username,
    u.email,
    u.api_token,
    u.is_active as user_active,
    u.created_at as user_created_at,
    t.id as tariff_id,
    t.name as tariff_name,
    t.description as tariff_description,
    t.minute_limit,
    t.day_limit,
    t.is_active as tariff_active
FROM users u
LEFT JOIN tariffs t ON u.tariff_id = t.id;

-- Создание представления для статистики использования тарифов
CREATE OR REPLACE VIEW tariff_usage_stats AS
SELECT 
    t.id as tariff_id,
    t.name as tariff_name,
    COUNT(u.id) as user_count,
    t.minute_limit,
    t.day_limit,
    AVG(arl.minute_count) as avg_minute_usage,
    AVG(arl.day_count) as avg_day_usage,
    MAX(arl.last_request_time) as last_activity
FROM tariffs t
LEFT JOIN users u ON t.id = u.tariff_id AND u.is_active = TRUE
LEFT JOIN api_rate_limits arl ON u.id = arl.user_id
WHERE t.is_active = TRUE
GROUP BY t.id, t.name, t.minute_limit, t.day_limit;

-- Создание процедуры для обновления тарифа пользователя
DELIMITER //
CREATE PROCEDURE UpdateUserTariff(
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
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = p_user_id) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Пользователь не найден';
    END IF;
    
    -- Проверяем существование тарифа
    IF NOT EXISTS (SELECT 1 FROM tariffs WHERE id = p_tariff_id AND is_active = TRUE) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Тариф не найден или неактивен';
    END IF;
    
    -- Обновляем тариф пользователя
    UPDATE users 
    SET tariff_id = p_tariff_id, updated_at = CURRENT_TIMESTAMP 
    WHERE id = p_user_id;
    
    -- Сбрасываем счетчики rate limiting для пользователя
    UPDATE api_rate_limits 
    SET minute_count = 0, day_count = 0, updated_at = CURRENT_TIMESTAMP 
    WHERE user_id = p_user_id;
    
    COMMIT;
END //
DELIMITER ;

-- Создание процедуры для получения лимитов пользователя
DELIMITER //
CREATE PROCEDURE GetUserLimits(
    IN p_user_id INT,
    OUT p_minute_limit INT,
    OUT p_day_limit INT
)
BEGIN
    SELECT t.minute_limit, t.day_limit
    INTO p_minute_limit, p_day_limit
    FROM users u
    JOIN tariffs t ON u.tariff_id = t.id
    WHERE u.id = p_user_id AND u.is_active = TRUE AND t.is_active = TRUE;
    
    -- Если пользователь не найден, используем дефолтные лимиты
    IF p_minute_limit IS NULL THEN
        SET p_minute_limit = 1000;
        SET p_day_limit = 10000;
    END IF;
END //
DELIMITER ;

-- Создание триггера для автоматического обновления updated_at
DELIMITER //
CREATE TRIGGER tariffs_updated_at 
    BEFORE UPDATE ON tariffs 
    FOR EACH ROW 
BEGIN
    SET NEW.updated_at = CURRENT_TIMESTAMP;
END //
DELIMITER ;

DELIMITER //
CREATE TRIGGER users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
BEGIN
    SET NEW.updated_at = CURRENT_TIMESTAMP;
END //
DELIMITER ;

-- Вывод информации о созданных таблицах
SELECT 'Таблица tariffs создана' as status;
SELECT COUNT(*) as tariff_count FROM tariffs;
SELECT 'Таблица users создана' as status;
SELECT COUNT(*) as user_count FROM users;
SELECT 'Таблица api_rate_limits создана' as status;
SELECT COUNT(*) as rate_limit_count FROM api_rate_limits;
