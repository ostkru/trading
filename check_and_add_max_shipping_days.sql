-- Проверка и добавление поля max_shipping_days в таблицы offers и orders

-- Проверяем, есть ли поле в таблице offers
SELECT COUNT(*) as column_exists 
FROM information_schema.columns 
WHERE table_schema = 'portaldata' 
  AND table_name = 'offers' 
  AND column_name = 'max_shipping_days';

-- Если поле не существует, добавляем его
ALTER TABLE offers ADD COLUMN IF NOT EXISTS max_shipping_days INT DEFAULT 7;

-- Проверяем, есть ли поле в таблице orders
SELECT COUNT(*) as column_exists 
FROM information_schema.columns 
WHERE table_schema = 'portaldata' 
  AND table_name = 'orders' 
  AND column_name = 'max_shipping_days';

-- Если поле не существует, добавляем его
ALTER TABLE orders ADD COLUMN IF NOT EXISTS max_shipping_days INT DEFAULT 7; 