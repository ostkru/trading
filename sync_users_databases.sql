-- Скрипт синхронизации таблиц users между базами pdata и portaldata

-- 1. Обновляем структуру таблицы users в portaldata для соответствия pdata
USE portaldata;

-- Добавляем недостающие поля в таблицу users
ALTER TABLE users 
ADD COLUMN username VARCHAR(50) UNIQUE AFTER id,
ADD COLUMN password_hash VARCHAR(255) AFTER email,
ADD COLUMN token_expires_at DATETIME AFTER api_token,
ADD COLUMN is_active TINYINT(1) DEFAULT 1 AFTER token_expires_at,
ADD COLUMN request_count INT DEFAULT 0 AFTER is_active,
ADD COLUMN last_request_at TIMESTAMP NULL AFTER request_count,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER created_at,
ADD COLUMN daily_request_count INT DEFAULT 0 AFTER updated_at,
ADD COLUMN last_daily_reset_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER daily_request_count,
ADD COLUMN minute_window_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER last_daily_reset_at,
ADD COLUMN day_window_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER minute_window_start,
ADD COLUMN webhook_url VARCHAR(255) AFTER day_window_start;

-- Создаем индексы для производительности
CREATE INDEX idx_users_token_expires_at ON users(token_expires_at);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_request_count ON users(request_count);

-- 2. Создаем процедуру для синхронизации пользователей из pdata в portaldata
DELIMITER $$

CREATE PROCEDURE SyncUsersFromPdata()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE pdata_id INT;
    DECLARE pdata_username VARCHAR(50);
    DECLARE pdata_email VARCHAR(100);
    DECLARE pdata_password_hash VARCHAR(255);
    DECLARE pdata_api_token VARCHAR(64);
    DECLARE pdata_token_expires_at DATETIME;
    DECLARE pdata_is_active TINYINT(1);
    DECLARE pdata_request_count INT;
    DECLARE pdata_last_request_at TIMESTAMP;
    DECLARE pdata_created_at TIMESTAMP;
    DECLARE pdata_updated_at TIMESTAMP;
    DECLARE pdata_daily_request_count INT;
    DECLARE pdata_last_daily_reset_at TIMESTAMP;
    DECLARE pdata_minute_window_start TIMESTAMP;
    DECLARE pdata_day_window_start TIMESTAMP;
    DECLARE pdata_webhook_url VARCHAR(255);
    
    DECLARE cur CURSOR FOR 
        SELECT id, username, email, password_hash, api_token, token_expires_at, 
               is_active, request_count, last_request_at, created_at, updated_at,
               daily_request_count, last_daily_reset_at, minute_window_start, 
               day_window_start, webhook_url
        FROM pdata.users;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    OPEN cur;
    
    read_loop: LOOP
        FETCH cur INTO pdata_id, pdata_username, pdata_email, pdata_password_hash, 
                       pdata_api_token, pdata_token_expires_at, pdata_is_active, 
                       pdata_request_count, pdata_last_request_at, pdata_created_at, 
                       pdata_updated_at, pdata_daily_request_count, pdata_last_daily_reset_at,
                       pdata_minute_window_start, pdata_day_window_start, pdata_webhook_url;
        
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- Вставляем или обновляем пользователя в portaldata
        INSERT INTO portaldata.users (
            id, username, name, email, password, password_hash, api_token, 
            token_expires_at, is_active, request_count, last_request_at, 
            created_at, updated_at, daily_request_count, last_daily_reset_at,
            minute_window_start, day_window_start, webhook_url
        ) VALUES (
            pdata_id, pdata_username, pdata_username, pdata_email, '', pdata_password_hash, 
            pdata_api_token, pdata_token_expires_at, pdata_is_active, pdata_request_count, 
            pdata_last_request_at, pdata_created_at, pdata_updated_at, pdata_daily_request_count,
            pdata_last_daily_reset_at, pdata_minute_window_start, pdata_day_window_start, 
            pdata_webhook_url
        ) ON DUPLICATE KEY UPDATE
            username = VALUES(username),
            name = VALUES(name),
            email = VALUES(email),
            password_hash = VALUES(password_hash),
            api_token = VALUES(api_token),
            token_expires_at = VALUES(token_expires_at),
            is_active = VALUES(is_active),
            request_count = VALUES(request_count),
            last_request_at = VALUES(last_request_at),
            updated_at = VALUES(updated_at),
            daily_request_count = VALUES(daily_request_count),
            last_daily_reset_at = VALUES(last_daily_reset_at),
            minute_window_start = VALUES(minute_window_start),
            day_window_start = VALUES(day_window_start),
            webhook_url = VALUES(webhook_url);
    END LOOP;
    
    CLOSE cur;
END$$

DELIMITER ;

-- 3. Создаем процедуру для синхронизации пользователей из portaldata в pdata
DELIMITER $$

CREATE PROCEDURE SyncUsersToPdata()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE portal_id INT;
    DECLARE portal_username VARCHAR(50);
    DECLARE portal_email VARCHAR(100);
    DECLARE portal_password_hash VARCHAR(255);
    DECLARE portal_api_token VARCHAR(64);
    DECLARE portal_token_expires_at DATETIME;
    DECLARE portal_is_active TINYINT(1);
    DECLARE portal_request_count INT;
    DECLARE portal_last_request_at TIMESTAMP;
    DECLARE portal_created_at TIMESTAMP;
    DECLARE portal_updated_at TIMESTAMP;
    DECLARE portal_daily_request_count INT;
    DECLARE portal_last_daily_reset_at TIMESTAMP;
    DECLARE portal_minute_window_start TIMESTAMP;
    DECLARE portal_day_window_start TIMESTAMP;
    DECLARE portal_webhook_url VARCHAR(255);
    
    DECLARE cur CURSOR FOR 
        SELECT id, username, email, password_hash, api_token, token_expires_at, 
               is_active, request_count, last_request_at, created_at, updated_at,
               daily_request_count, last_daily_reset_at, minute_window_start, 
               day_window_start, webhook_url
        FROM portaldata.users;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    OPEN cur;
    
    read_loop: LOOP
        FETCH cur INTO portal_id, portal_username, portal_email, portal_password_hash, 
                       portal_api_token, portal_token_expires_at, portal_is_active, 
                       portal_request_count, portal_last_request_at, portal_created_at, 
                       portal_updated_at, portal_daily_request_count, portal_last_daily_reset_at,
                       portal_minute_window_start, portal_day_window_start, portal_webhook_url;
        
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- Вставляем или обновляем пользователя в pdata
        INSERT INTO pdata.users (
            id, username, email, password_hash, api_token, token_expires_at, 
            is_active, request_count, last_request_at, created_at, updated_at,
            daily_request_count, last_daily_reset_at, minute_window_start, 
            day_window_start, webhook_url
        ) VALUES (
            portal_id, portal_username, portal_email, portal_password_hash, 
            portal_api_token, portal_token_expires_at, portal_is_active, 
            portal_request_count, portal_last_request_at, portal_created_at, 
            portal_updated_at, portal_daily_request_count, portal_last_daily_reset_at,
            portal_minute_window_start, portal_day_window_start, portal_webhook_url
        ) ON DUPLICATE KEY UPDATE
            username = VALUES(username),
            email = VALUES(email),
            password_hash = VALUES(password_hash),
            api_token = VALUES(api_token),
            token_expires_at = VALUES(token_expires_at),
            is_active = VALUES(is_active),
            request_count = VALUES(request_count),
            last_request_at = VALUES(last_request_at),
            updated_at = VALUES(updated_at),
            daily_request_count = VALUES(daily_request_count),
            last_daily_reset_at = VALUES(last_daily_reset_at),
            minute_window_start = VALUES(minute_window_start),
            day_window_start = VALUES(day_window_start),
            webhook_url = VALUES(webhook_url);
    END LOOP;
    
    CLOSE cur;
END$$

DELIMITER ;

-- 4. Создаем событие для автоматической синхронизации каждые 5 минут
CREATE EVENT IF NOT EXISTS SyncUsersEvent
ON SCHEDULE EVERY 5 MINUTE
DO
BEGIN
    CALL SyncUsersFromPdata();
    CALL SyncUsersToPdata();
END;

-- 5. Включаем планировщик событий
SET GLOBAL event_scheduler = ON;

-- 6. Запускаем первоначальную синхронизацию
CALL SyncUsersFromPdata();
CALL SyncUsersToPdata();

-- Выводим результат синхронизации
SELECT 'Синхронизация завершена' AS status;
SELECT COUNT(*) AS users_in_pdata FROM pdata.users;
SELECT COUNT(*) AS users_in_portaldata FROM portaldata.users; 