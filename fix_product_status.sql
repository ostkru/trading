-- Исправление статусов продуктов - все должны быть 'pending'
-- Выполнить: mysql -u root -p123456 portaldata < fix_product_status.sql

USE portaldata;

-- Обновляем все продукты на статус 'pending'
UPDATE products SET status = 'pending';

-- Проверяем результат
SELECT status, COUNT(*) as count FROM products GROUP BY status;

-- Показываем примеры записей с их category_id и brand_id
SELECT id, name, status, category_id, brand_id FROM products LIMIT 10; 