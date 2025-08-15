-- SQL скрипт для оптимизации производительности базы данных
-- Выполнить в базе данных portaldata

-- 1. Добавление составных индексов для ускорения запросов
ALTER TABLE products ADD INDEX idx_user_status (user_id, status);
ALTER TABLE products ADD INDEX idx_brand_category (brand, category);
ALTER TABLE products ADD INDEX idx_created_at (created_at);
ALTER TABLE products ADD INDEX idx_updated_at (updated_at);

-- 2. Оптимизация индекса для media таблицы
ALTER TABLE media ADD INDEX idx_product_created (product_id, created_at);

-- 3. Анализ и оптимизация таблиц
ANALYZE TABLE products;
ANALYZE TABLE media;

-- 4. Проверка текущих индексов
SHOW INDEX FROM products;
SHOW INDEX FROM media;

-- 5. Статистика по таблицам
SELECT 
    table_name,
    table_rows,
    data_length,
    index_length,
    (data_length + index_length) as total_size
FROM information_schema.tables 
WHERE table_schema = 'portaldata' 
AND table_name IN ('products', 'media');

-- 6. Проверка размера InnoDB буфера
SHOW VARIABLES LIKE 'innodb_buffer_pool_size';
SHOW VARIABLES LIKE 'innodb_log_file_size';
SHOW VARIABLES LIKE 'innodb_flush_log_at_trx_commit';

-- 7. Рекомендуемые настройки для оптимизации (выполнить в my.cnf)
-- innodb_buffer_pool_size = 1G (минимум)
-- innodb_log_file_size = 256M
-- innodb_flush_log_at_trx_commit = 2 (для лучшей производительности)
-- innodb_flush_method = O_DIRECT
-- innodb_file_per_table = 1
