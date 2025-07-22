-- Финальный скрипт синхронизации пользователей между pdata и portaldata

USE portaldata;

-- Создаем процедуру для синхронизации из pdata в portaldata
DELIMITER $$

CREATE PROCEDURE SyncUsersFromPdata()
BEGIN
    INSERT INTO portaldata.users (
        id, username, name, email, password, password_hash, api_token, 
        token_expires_at, is_active, request_count, last_request_at, 
        created_at, updated_at, daily_request_count, last_daily_reset_at,
        minute_window_start, day_window_start, webhook_url
    )
    SELECT 
        id, username, username as name, email, '' as password, password_hash, api_token,
        token_expires_at, is_active, request_count, last_request_at,
        created_at, updated_at, daily_request_count, last_daily_reset_at,
        minute_window_start, day_window_start, webhook_url
    FROM pdata.users
    ON DUPLICATE KEY UPDATE
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
END$$

DELIMITER ;

-- Создаем процедуру для синхронизации из portaldata в pdata
DELIMITER $$

CREATE PROCEDURE SyncUsersToPdata()
BEGIN
    INSERT INTO pdata.users (
        id, username, email, password_hash, api_token, token_expires_at, 
        is_active, request_count, last_request_at, created_at, updated_at,
        daily_request_count, last_daily_reset_at, minute_window_start, 
        day_window_start, webhook_url
    )
    SELECT 
        id, username, email, password_hash, api_token, token_expires_at,
        is_active, request_count, last_request_at, created_at, updated_at,
        daily_request_count, last_daily_reset_at, minute_window_start,
        day_window_start, webhook_url
    FROM portaldata.users
    ON DUPLICATE KEY UPDATE
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
END$$

DELIMITER ;

-- Запускаем синхронизацию
CALL SyncUsersFromPdata();
CALL SyncUsersToPdata();

-- Показываем результаты
SELECT 'Синхронизация завершена' AS status;
SELECT COUNT(*) AS users_in_pdata FROM pdata.users;
SELECT COUNT(*) AS users_in_portaldata FROM portaldata.users; 