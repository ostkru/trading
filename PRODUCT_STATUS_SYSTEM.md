# Система статусов продуктов

## Обзор

Система статусов продуктов реализована для контроля процесса классификации товаров. Все продукты имеют статус `pending` и классифицируются внешним скриптом на основе наличия `category_id` и `brand_id`.

## Статусы продуктов

### `pending` (по умолчанию)
- **Описание**: Все продукты создаются со статусом `pending`
- **Изменение**: Статус не может быть изменен через API
- **Использование**: Внешний скрипт классификации изменяет `category_id` и `brand_id`

## Фильтры для просмотра продуктов

### `my` (по умолчанию)
- Показывает только продукты текущего пользователя
- **SQL**: `WHERE user_id = ?`

### `others`
- Показывает продукты других пользователей
- **SQL**: `WHERE user_id != ?`

### `pending`
- Показывает все продукты со статусом `pending` (ожидающие классификации)
- **SQL**: `WHERE status = 'pending'`

### `not_classified`
- Показывает продукты со статусом `pending`, у которых отсутствует `category_id` или `brand_id`
- **SQL**: `WHERE status = 'pending' AND (category_id IS NULL OR brand_id IS NULL)`

### `classified`
- Показывает продукты со статусом `pending`, у которых есть и `category_id`, и `brand_id`
- **SQL**: `WHERE status = 'pending' AND category_id IS NOT NULL AND brand_id IS NOT NULL`

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

# Продукты, требующие классификации
curl "http://localhost:8095/api/v1/products?owner=not_classified&api_key=YOUR_API_KEY"

# Классифицированные продукты
curl "http://localhost:8095/api/v1/products?owner=classified&api_key=YOUR_API_KEY"
```

### POST /api/v1/products
**Создание продукта:**
- Статус автоматически устанавливается в `pending`
- `category_id` и `brand_id` изначально `NULL`

### PUT /api/v1/products/:id
**Обновление продукта:**
- Поле `status` недоступно для изменения через API
- Можно изменять только: name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode

## Внешний скрипт классификации

### Обязательные действия:
1. **Найти продукты со статусом `pending`**
2. **Определить `category_id` и `brand_id`**
3. **Обновить поля в базе данных**

### Пример SQL для внешнего скрипта:
```sql
-- Найти продукты для классификации
SELECT id, name, brand, category FROM products 
WHERE status = 'pending' AND (category_id IS NULL OR brand_id IS NULL);

-- Обновить классифицированные продукты
UPDATE products 
SET category_id = ?, brand_id = ? 
WHERE id = ?;
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
    status ENUM('pending', 'classified', 'not_classified') DEFAULT 'pending',
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

# Проверить количество неклассифицированных продуктов
curl "http://localhost:8095/api/v1/products?owner=not_classified&api_key=YOUR_API_KEY" | jq '.total'

# Проверить количество классифицированных продуктов
curl "http://localhost:8095/api/v1/products?owner=classified&api_key=YOUR_API_KEY" | jq '.total'
```

## Безопасность

- Статус продуктов не может быть изменен через API
- Только внешний скрипт может обновлять `category_id` и `brand_id`
- Офферы создаются только для классифицированных продуктов
- Все изменения логируются в `app.log` 