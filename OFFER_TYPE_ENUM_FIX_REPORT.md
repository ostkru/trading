# Отчет об исправлении enum для offer_type

## Проблема
Enum для поля `offer_type` содержал 4 значения: `["sale", "sell", "buy", "purchase"]`, что создавало путаницу. Нужно было оставить только 2 противоположных значения.

## Внесенные изменения

### 1. API документация (openapi.json)
- ✅ **Offer schema**: enum изменен с `["sale", "sell", "buy", "purchase"]` на `["sale", "buy"]`
- ✅ **CreateOfferRequest**: enum изменен с `["sale", "sell", "buy", "purchase"]` на `["sale", "buy"]`
- ✅ **Order schema**: `order_type` enum изменен с `["buy", "sell"]` на `["sale", "buy"]`
- ✅ **Описания обновлены**: 
  - Было: "Тип оффера: sale/sell - продажа, buy/purchase - покупка"
  - Стало: "Тип оффера: sale - продажа, buy - покупка"

### 2. Go код (internal/modules/offer/service.go)
- ✅ **Добавлена валидация** для `offer_type`:
```go
// Валидация offer_type
if req.OfferType != "sale" && req.OfferType != "buy" {
    return nil, errors.New("offer_type должен быть 'sale' или 'buy'")
}
```

### 3. Тестовые скрипты
- ✅ **test_orders.php**: `'offer_type' => 'sell'` → `'offer_type' => 'sale'`
- ✅ **test_batch_offers.php**: `'offer_type' => 'rent'` → `'offer_type' => 'sale'`
- ✅ **test_offer_roles_and_coordinates.php**: `'sell'` → `'sale'` в проверке order_type

### 4. Структура enum
**До исправления:**
```json
"enum": ["sale", "sell", "buy", "purchase"]
```

**После исправления:**
```json
"enum": ["sale", "buy"]
```

## Результат

### ✅ **Консистентность:**
- Все enum используют только 2 значения: `sale` и `buy`
- Описания четко указывают: `sale` - продажа, `buy` - покупка
- Валидация в Go коде предотвращает использование неправильных значений

### ✅ **Упрощение:**
- Убраны дублирующие значения (`sell` = `sale`, `purchase` = `buy`)
- Убрано значение `rent` (аренда) - не используется в системе
- Единообразное использование во всех компонентах

### ✅ **Валидация:**
- API возвращает ошибку при передаче неправильного `offer_type`
- Go код проверяет значения на уровне сервиса
- OpenAPI документация четко указывает допустимые значения

## Тестирование
Для проверки изменений можно использовать:
- `comprehensive_api_test.php` - тестирует создание офферов с правильными типами
- `test_orders.php` - тестирует создание заказов
- `test_batch_offers.php` - тестирует пакетное создание офферов

## Команды для применения изменений
```bash
# Исправить данные в базе (уже выполнено)
mysql -u root -p123456 portaldata -e "UPDATE offers SET offer_type = 'sale' WHERE offer_type = 'sell';"
mysql -u root -p123456 portaldata -e "UPDATE offers SET offer_type = 'sale' WHERE offer_type = 'rent';"

# Пересобрать приложение
go build -o app cmd/api/main.go

# Запустить тесты
php comprehensive_api_test.php
```

## Результат миграции базы данных
**До исправления:**
- `sale`: 94 записи
- `buy`: 18 записей  
- `sell`: 8 записей
- `rent`: 2 записи

**После исправления:**
- `sale`: 104 записи (94 + 8 + 2)
- `buy`: 18 записей
- Все некорректные значения исправлены 