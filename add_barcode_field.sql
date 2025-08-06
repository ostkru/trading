-- Добавление поля barcode в таблицу products
-- Выполнить: mysql -u root -p123456 portaldata < add_barcode_field.sql

USE portaldata;

-- Добавляем поле barcode как VARCHAR(50) с возможностью NULL
ALTER TABLE products ADD COLUMN barcode VARCHAR(50) NULL COMMENT 'Штрих-код продукта';

-- Создаем индекс для быстрого поиска по штрих-коду
CREATE INDEX idx_products_barcode ON products(barcode);

-- Проверяем, что поле добавлено
DESCRIBE products;

-- Показываем примеры записей с новым полем
SELECT id, name, vendor_article, barcode FROM products LIMIT 5; 