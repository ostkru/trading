-- Обновление статусов продуктов для поддержки классификации
-- Выполнить: mysql -u root -p123456 portaldata < update_product_status.sql

USE portaldata;

-- Обновляем enum для статусов продуктов
ALTER TABLE products MODIFY COLUMN status ENUM('pending', 'classified', 'not_classified') DEFAULT 'pending';

-- Обновляем существующие записи
-- Если есть category_id и brand_id, устанавливаем статус classified
UPDATE products SET status = 'classified' WHERE category_id IS NOT NULL AND brand_id IS NOT NULL;

-- Если нет category_id или brand_id, устанавливаем статус not_classified
UPDATE products SET status = 'not_classified' WHERE (category_id IS NULL OR brand_id IS NULL) AND status != 'classified';

-- Проверяем результат
SELECT status, COUNT(*) as count FROM products GROUP BY status;

-- Показываем примеры записей
SELECT id, name, status, category_id, brand_id FROM products LIMIT 10; 