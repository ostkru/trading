-- Добавление колонок category_id и brand_id в таблицу products
-- Эти колонки будут использоваться для автоматического определения категорий и брендов

-- Добавляем колонку category_id
ALTER TABLE products ADD COLUMN category_id INT NULL;

-- Добавляем колонку brand_id  
ALTER TABLE products ADD COLUMN brand_id INT NULL;

-- Проверяем, что колонки добавлены
DESCRIBE products; 