USE portaldata;

SET GLOBAL event_scheduler = ON;

DELIMITER $$

CREATE EVENT IF NOT EXISTS SyncUsersEvent
ON SCHEDULE EVERY 5 MINUTE
DO
BEGIN
    CALL SyncUsersFromPdata();
    CALL SyncUsersToPdata();
END$$

DELIMITER ;

SHOW EVENTS; 