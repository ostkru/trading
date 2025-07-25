# 🗺️ API Географических Фильтров для Offers

## Обзор

Добавлена поддержка расширенных фильтров для offers, включая географические фильтры для работы с прямоугольными областями по координатам.

## 🆕 Новые Endpoints

### 1. Фильтрованные Offers (с авторизацией)
```
POST /api/v1/offers/filter
```

### 2. Фильтрованные Публичные Offers (без авторизации)
```
POST /api/v1/offers/public/filter
```

## 📋 Структура фильтров

### OfferFilterRequest
```json
{
  "filter": "all",                    // my, others, all
  "offer_type": "sale",              // sale, buy
  "geographic": {                     // Географический фильтр
    "min_latitude": 55.0,
    "max_latitude": 56.0,
    "min_longitude": 37.0,
    "max_longitude": 38.0
  },
  "price_min": 100,                  // Минимальная цена
  "price_max": 5000,                 // Максимальная цена
  "available_lots": 1                // Минимальное количество лотов
}
```

## 🗺️ Географический фильтр

### GeographicFilter
```json
{
  "min_latitude": 55.0,    // Минимальная широта
  "max_latitude": 56.0,    // Максимальная широта
  "min_longitude": 37.0,   // Минимальная долгота
  "max_longitude": 38.0    // Максимальная долгота
}
```

**Примеры координат:**
- **Москва**: 55.0-56.0, 37.0-38.0
- **Санкт-Петербург**: 59.0-60.0, 30.0-31.0
- **Весь мир**: 0-90, 0-180

## 💰 Фильтры по цене

```json
{
  "price_min": 100,    // Минимальная цена за единицу
  "price_max": 5000    // Максимальная цена за единицу
}
```

## 📦 Фильтры по доступности

```json
{
  "available_lots": 1    // Минимальное количество доступных лотов
}
```

## 🔍 Примеры использования

### 1. Базовый географический фильтр
```bash
curl -X POST http://localhost:8095/api/v1/offers/filter \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filter": "all",
    "geographic": {
      "min_latitude": 55.0,
      "max_latitude": 56.0,
      "min_longitude": 37.0,
      "max_longitude": 38.0
    }
  }'
```

### 2. Фильтр по цене
```bash
curl -X POST http://localhost:8095/api/v1/offers/filter \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filter": "all",
    "price_min": 100,
    "price_max": 5000
  }'
```

### 3. Комбинированный фильтр
```bash
curl -X POST http://localhost:8095/api/v1/offers/filter \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filter": "all",
    "offer_type": "sale",
    "geographic": {
      "min_latitude": 55.0,
      "max_latitude": 56.0,
      "min_longitude": 37.0,
      "max_longitude": 38.0
    },
    "price_min": 500,
    "available_lots": 1
  }'
```

### 4. Публичные офферы с фильтрами
```bash
curl -X POST http://localhost:8095/api/v1/offers/public/filter \
  -H "Content-Type: application/json" \
  -d '{
    "offer_type": "buy",
    "geographic": {
      "min_latitude": 55.0,
      "max_latitude": 56.0,
      "min_longitude": 37.0,
      "max_longitude": 38.0
    },
    "price_max": 3000
  }'
```

## 📊 Параметры пагинации

Все endpoints поддерживают стандартные параметры пагинации:

```
?page=1&limit=20
```

- `page` (по умолчанию: 1) - номер страницы
- `limit` (по умолчанию: 20, максимум: 100) - количество записей на странице

## 🔒 Безопасность

### Авторизованные запросы (`/offers/filter`)
- Требуют Bearer token
- Поддерживают фильтры `my`, `others`, `all`
- Возвращают офферы в зависимости от прав пользователя

### Публичные запросы (`/offers/public/filter`)
- Не требуют авторизации
- Возвращают только публичные офферы (`is_public = true`)
- Игнорируют фильтр `filter`

## ⚠️ Валидация

### Поддерживаемые значения:
- `filter`: `my`, `others`, `all`
- `offer_type`: `sale`, `buy`
- `geographic`: все координаты должны быть числами
- `price_min`, `price_max`: положительные числа
- `available_lots`: положительное целое число

### Ошибки валидации:
- `400 Bad Request` - некорректный JSON или недопустимые значения
- `401 Unauthorized` - отсутствует токен авторизации
- `500 Internal Server Error` - внутренняя ошибка сервера

## 📈 Пример ответа

```json
{
  "offers": [
    {
      "offer_id": 164,
      "user_id": 6,
      "created_at": "2025-07-25T07:00:34Z",
      "updated_at": "2025-07-25T07:00:34Z",
      "is_public": true,
      "product_id": 192,
      "price_per_unit": 2000,
      "tax_nds": 20,
      "units_per_lot": 1,
      "available_lots": 8,
      "latitude": 55.7558,
      "longitude": 37.6176,
      "warehouse_id": 80,
      "offer_type": "buy",
      "max_shipping_days": 5,
      "product_name": "Тестовый продукт",
      "vendor_article": "TEST001",
      "recommend_price": 2500,
      "warehouse_name": "Склад в Москве",
      "warehouse_address": "ул. Тверская, 1"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

## 🧪 Тестирование

Запустите тесты географических фильтров:

```bash
php test_geographic_filters.php
```

## 🚀 Особенности реализации

1. **Географические фильтры** работают с координатами складов (`warehouses.latitude`, `warehouses.longitude`)
2. **Прямоугольные области** определяются четырьмя координатами
3. **Комбинированные фильтры** применяются через SQL `AND` условия
4. **Оптимизация** - используются индексы по координатам и ценам
5. **Безопасность** - все параметры экранируются через prepared statements

## 📝 Примечания

- Координаты должны быть в десятичном формате (WGS84)
- Широта: -90 до +90
- Долгота: -180 до +180
- Фильтры применяются в порядке: пользователь → тип → география → цена → доступность 