-- Исправление значений offer_type в базе данных
-- Приводим к единому стандарту: sale и buy

-- Заменяем 'sell' на 'sale' (продажа)
UPDATE offers SET offer_type = 'sale' WHERE offer_type = 'sell';

-- Заменяем 'rent' на 'sale' (аренда = продажа)
UPDATE offers SET offer_type = 'sale' WHERE offer_type = 'rent';

-- Проверяем результат
SELECT DISTINCT offer_type, COUNT(*) as count FROM offers GROUP BY offer_type ORDER BY count DESC;

-- Проверяем, что все значения корректны
SELECT COUNT(*) as total_offers FROM offers WHERE offer_type NOT IN ('sale', 'buy'); 