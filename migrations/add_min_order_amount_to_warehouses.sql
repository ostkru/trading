-- Добавление колонки min_order_amount в таблицу warehouses
-- Дата: 2025-09-22

-- Добавляем колонку min_order_amount с значением по умолчанию 0
ALTER TABLE warehouses 
ADD COLUMN min_order_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 
COMMENT 'Минимальная сумма заказа для склада';

-- Оставляем существующие записи с значением 0 (без минимальной суммы заказа)
-- UPDATE warehouses 
-- SET min_order_amount = 0.00 
-- WHERE min_order_amount = 0.00;

-- Добавляем индекс для оптимизации запросов по минимальной сумме заказа
CREATE INDEX idx_warehouses_min_order_amount ON warehouses(min_order_amount);

-- Добавляем комментарий к таблице
ALTER TABLE warehouses COMMENT = 'Склады с информацией о минимальной сумме заказа';
