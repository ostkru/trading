# Система статусов продуктов

## Обзор

Система статусов продуктов реализована для контроля процесса классификации товаров. Все продукты создаются со статусом `processing` и автоматически классифицируются асинхронным воркером на основе анализа названия, бренда и категории.

## Статусы продуктов

### `processing` (по умолчанию)
- **Описание**: Все продукты создаются со статусом `processing`
- **Изменение**: Статус автоматически изменяется воркером классификации
- **Использование**: Асинхронный воркер анализирует продукт и обновляет `category_id` и `brand_id`

## Фильтры для просмотра продуктов

### `my` (по умолчанию)
- Показывает только продукты текущего пользователя
- **SQL**: `WHERE user_id = ?`

### `others`
- Показывает продукты других пользователей
- **SQL**: `WHERE user_id != ?`

### `processing`
- Показывает все продукты со статусом `processing` (ожидающие классификации)
- **SQL**: `WHERE status = 'processing'`

### `not_classified`
- Показывает продукты со статусом `not_classified` (не удалось классифицировать)
- **SQL**: `WHERE status = 'not_classified'`

### `classified`
- Показывает продукты со статусом `classified` (успешно классифицированы)
- **SQL**: `WHERE status = 'classified'`

## Ограничения для офферов

### Проверка при создании оффера
- Офферы могут создаваться только для продуктов с заполненными `category_id` и `brand_id`
- **Ошибка**: "Нельзя создать оффер для продукта без category_id или brand_id. Продукт должен быть классифицирован"

### Логика проверки
```sql
SELECT category_id, brand_id FROM products WHERE id = ?
```

## API Endpoints

### GET /api/v1/products
**Параметры:**
- `owner` - фильтр (my, others, pending, not_classified, classified)
- `page` - номер страницы
- `limit` - количество записей на странице

**Примеры:**
```bash
# Мои продукты
curl "http://localhost:8095/api/v1/products?owner=my&api_key=YOUR_API_KEY"

# Продукты в обработке
curl "http://localhost:8095/api/v1/products?owner=processing&api_key=YOUR_API_KEY"

# Продукты, требующие классификации
curl "http://localhost:8095/api/v1/products?owner=not_classified&api_key=YOUR_API_KEY"

# Классифицированные продукты
curl "http://localhost:8095/api/v1/products?owner=classified&api_key=YOUR_API_KEY"
```

### POST /api/v1/products
**Создание продукта:**
- Статус автоматически устанавливается в `processing`
- `category_id` и `brand_id` изначально `NULL`
- Запускается асинхронная классификация

### PUT /api/v1/products/:id
**Обновление продукта:**
- Поле `status` недоступно для изменения через API
- Можно изменять только: name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode

## Асинхронная классификация

### Автоматический процесс:
1. **При создании продукта** запускается асинхронная классификация
2. **Воркер классификации** анализирует название, бренд и категорию
3. **Требования для `classified`:**
   - `confidence.category >= 0.99` → обновляется `category_id`
   - `confidence.brand >= 0.99` → обновляется `brand_id`
   - **Минимум один ID** должен быть найден с вероятностью ≥ 0.99

### Результаты классификации:

**✅ Успешная классификация:**
- **Статус:** `classified`
- **Поля:** `category_id` и/или `brand_id` заполнены
- **Условие:** Найден хотя бы один ID с вероятностью ≥ 0.99

**❌ Неудачная классификация:**
- **Статус:** `not_classified`
- **Поля:** `category_id = NULL`, `brand_id = NULL`
- **Условие:** Не найдены ID с достаточной вероятностью

### Пример SQL для мониторинга:
```sql
-- Найти продукты для классификации
SELECT id, name, brand, category FROM products 
WHERE status = 'processing';

-- Найти классифицированные продукты
SELECT id, name, category_id, brand_id FROM products 
WHERE status = 'classified';

-- Найти неклассифицированные продукты
SELECT id, name FROM products 
WHERE status = 'not_classified';
```

## Структура базы данных

### Таблица `products`
```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    vendor_article VARCHAR(255) NOT NULL,
    recommend_price DECIMAL(10,2),
    brand VARCHAR(255),
    category VARCHAR(255),
    brand_id INT NULL,
    category_id INT NULL,
    description TEXT,
    barcode VARCHAR(50) NULL,
    status ENUM('processing', 'classified', 'not_classified') DEFAULT 'processing',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id INT
);
```

## Тестирование

### Запуск тестов
```bash
./test_product_status.sh
```

### Проверка статусов
```bash
# Проверить статус продукта
curl "http://localhost:8095/api/v1/products/1?api_key=YOUR_API_KEY" | jq '.status'

# Проверить количество продуктов в обработке
curl "http://localhost:8095/api/v1/products?owner=processing&api_key=YOUR_API_KEY" | jq '.total'

# Проверить количество неклассифицированных продуктов
curl "http://localhost:8095/api/v1/products?owner=not_classified&api_key=YOUR_API_KEY" | jq '.total'

# Проверить количество классифицированных продуктов
curl "http://localhost:8095/api/v1/products?owner=classified&api_key=YOUR_API_KEY" | jq '.total'
```

## Безопасность

- Статус продуктов не может быть изменен через API
- Только асинхронный воркер классификации может обновлять `category_id` и `brand_id`
- Офферы создаются только для классифицированных продуктов
- Все изменения логируются в `app.log`

## Мониторинг классификации

### Логи воркера классификации:
```
🚀 Воркер классификации запущен
✅ Задача классификации добавлена в очередь для продукта 123
🔍 Обрабатываю классификацию продукта 123: Телевизор Samsung
📊 Результат классификации продукта 123: категория=15 (точность=0.99), бренд=42 (точность=0.98)
🎯 Продукт 123 успешно классифицирован
✅ Продукт 123 обновлен: category_id=15, brand_id=42, status=classified
```

### Статистика по статусам:
```sql
-- Общая статистика
SELECT status, COUNT(*) as count FROM products GROUP BY status;

-- Продукты в обработке
SELECT COUNT(*) FROM products WHERE status = 'processing';

-- Успешно классифицированные
SELECT COUNT(*) FROM products WHERE status = 'classified';

-- Неклассифицированные
SELECT COUNT(*) FROM products WHERE status = 'not_classified';
``` 