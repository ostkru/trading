-- Добавление поля max_shipping_days в таблицу orders
ALTER TABLE orders ADD COLUMN max_shipping_days INT DEFAULT 7; 