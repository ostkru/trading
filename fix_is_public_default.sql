-- Исправление значений по умолчанию для поля is_public
-- Обновляем все NULL значения на true

-- Для MySQL
UPDATE offers SET is_public = 1 WHERE is_public IS NULL;

-- Проверяем результат
SELECT offer_id, is_public FROM offers WHERE is_public IS NULL; 