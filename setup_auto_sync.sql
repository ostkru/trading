-- Настройка автоматической синхронизации пользователей

USE portaldata;

-- Включаем планировщик событий
SET GLOBAL event_scheduler = ON;

-- Создаем событие для автоматической синхронизации каждые 5 минут
CREATE EVENT IF NOT EXISTS SyncUsersEvent
ON SCHEDULE EVERY 5 MINUTE
DO
BEGIN
    CALL SyncUsersFromPdata();
    CALL SyncUsersToPdata();
END;

-- Проверяем, что событие создано
SHOW EVENTS;

-- Показываем статус планировщика
SHOW VARIABLES LIKE 'event_scheduler'; 