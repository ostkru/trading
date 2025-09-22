# 🎯 Интеграция тарифов с Rate Limiting - Руководство

## 🔍 Анализ текущей ситуации

После анализа проекта выяснилось, что:

1. **Таблица `tariffs` уже существует** и содержит тарифы с лимитами в JSON поле `features`
2. **Таблица `users` НЕ имеет связи с тарифами** - отсутствует поле `tariff_id`
3. **Rate limiting использует фиксированные лимиты** (60/мин, 1000/день) вместо тарифов
4. **Лимиты хранятся в JSON** как `daily_requests_limit`

## 📊 Текущие тарифы в системе

| ID | Название | Цена | Дневной лимит | Описание |
|----|----------|------|---------------|----------|
| 1 | Старт | 0.00 ₽ | 1,000 | Базовый тариф для начала работы |
| 2 | Продвинутый | 10,000 ₽ | 2,000 | Для активных пользователей |
| 3 | Мастер | 16,000 ₽ | 4,000 | Профессиональный тариф для торговли |
| 4 | Профи | 24,000 ₽ | 8,000 | Для торговли и высокоуровневой аналитики |
| 5 | Эксперт | 36,000 ₽ | 15,000 | Максимальные возможности |

## 🚀 Решение: Интеграция существующих тарифов

### Что было сделано:

1. **Добавлено поле `tariff_id`** в таблицу `users`
2. **Созданы функции** для получения лимитов из тарифов:
   - `GetUserDailyLimit(user_id)` - дневной лимит
   - `GetUserMinuteLimit(user_id)` - минутный лимит (рассчитывается автоматически)
3. **Обновлен сервис rate limiting** для использования тарифов
4. **Созданы представления** для удобного просмотра данных
5. **Добавлены процедуры** для управления тарифами

## 🛠️ Установка и настройка

### 1. Применение интеграции
```bash
cd /var/www/trading
./scripts/apply-tariff-integration.sh
```

### 2. Проверка результатов
```bash
# Проверка структуры таблицы users
mysql -u root -p -e "USE portaldata; DESCRIBE users;"

# Проверка функций
mysql -u root -p -e "USE portaldata; SHOW FUNCTION STATUS WHERE Name LIKE 'GetUser%';"

# Проверка представлений
mysql -u root -p -e "USE portaldata; SHOW TABLES LIKE '%tariff%';"
```

## 🔧 Как это работает

### Алгоритм получения лимитов:

1. **Пользователь делает запрос** с API ключом
2. **Система определяет userID** по API ключу
3. **Вызывается функция `GetUserDailyLimit(userID)`**:
   ```sql
   SELECT JSON_UNQUOTE(JSON_EXTRACT(features, '$.daily_requests_limit'))
   FROM users u
   JOIN tariffs t ON u.tariff_id = t.id
   WHERE u.id = userID AND u.is_active = 1 AND t.is_active = 1
   ```
4. **Вызывается функция `GetUserMinuteLimit(userID)`**:
   ```sql
   -- Минутный лимит = дневной лимит / 1440 (минут в дне)
   -- Но минимум 60 запросов в минуту
   ```
5. **Rate limiting проверяет** текущие счетчики против полученных лимитов

### Формула лимитов:
- **Дневной лимит**: берется из `features.daily_requests_limit` тарифа
- **Минутный лимит**: `MAX(60, daily_limit / 1440)`

## 📊 Примеры использования

### Просмотр пользователей с тарифами
```sql
USE portaldata;
SELECT * FROM user_tariffs_with_limits;
```

### Изменение тарифа пользователя
```sql
USE portaldata;
CALL ChangeUserTariff(1, 3); -- Пользователь 1 получает тариф "Мастер"
```

### Получение лимитов пользователя
```sql
USE portaldata;
CALL GetUserLimits(1, @minute_limit, @daily_limit);
SELECT @minute_limit, @daily_limit;
```

### Статистика использования тарифов
```sql
USE portaldata;
SELECT * FROM tariff_usage_statistics;
```

## 🔄 Обновление тарифов

### Изменение лимитов существующего тарифа
```sql
USE portaldata;
UPDATE tariffs 
SET features = JSON_SET(features, '$.daily_requests_limit', 5000)
WHERE id = 2;
```

### Добавление нового тарифа
```sql
USE portaldata;
INSERT INTO tariffs (name, price, duration_days, features) VALUES
('VIP', 50000.00, 365, '{"daily_requests_limit": 20000, "description": "VIP тариф", "features": ["Все возможности Эксперт", "Персональная поддержка"]}');
```

## 🧪 Тестирование

### Тест 1: Проверка лимитов разных пользователей
```bash
# Пользователь с тарифом "Старт" (1000/день)
curl -H "X-API-Key: test_api_key_1" http://localhost:8080/api/products

# Пользователь с тарифом "Эксперт" (15000/день)
curl -H "X-API-Key: test_api_key_5" http://localhost:8080/api/products
```

### Тест 2: Изменение тарифа и проверка лимитов
```sql
-- Меняем тариф пользователя
CALL ChangeUserTariff(1, 5);

-- Проверяем новые лимиты
SELECT GetUserDailyLimit(1), GetUserMinuteLimit(1);
```

## 📈 Мониторинг

### Представления для мониторинга:

#### `user_tariffs_with_limits`
Показывает всех пользователей с их тарифами и лимитами:
```sql
SELECT 
    user_id,
    username,
    tariff_name,
    daily_limit,
    minute_limit
FROM user_tariffs_with_limits
WHERE user_active = 1;
```

#### `tariff_usage_statistics`
Показывает статистику использования каждого тарифа:
```sql
SELECT 
    tariff_name,
    daily_limit,
    user_count,
    avg_daily_usage,
    last_activity
FROM tariff_usage_statistics
ORDER BY tariff_price;
```

## 🚨 Устранение неполадок

### Проблема: Пользователь не найден
```sql
-- Проверяем существование пользователя
SELECT id, username, is_active FROM users WHERE id = 1;

-- Проверяем назначенный тариф
SELECT u.id, u.username, t.name as tariff_name 
FROM users u 
LEFT JOIN tariffs t ON u.tariff_id = t.id 
WHERE u.id = 1;
```

### Проблема: Тариф не найден
```sql
-- Проверяем существование тарифа
SELECT id, name, is_active FROM tariffs WHERE id = 3;

-- Проверяем JSON структуру
SELECT id, name, features FROM tariffs WHERE id = 3;
```

### Проблема: Лимиты не применяются
```sql
-- Проверяем функции
SELECT GetUserDailyLimit(1), GetUserMinuteLimit(1);

-- Проверяем связь пользователь-тариф
SELECT u.id, u.tariff_id, t.name, t.features
FROM users u
JOIN tariffs t ON u.tariff_id = t.id
WHERE u.id = 1;
```

## 🔧 Конфигурация

### Переменные окружения:
```bash
# MySQL для тарифов
DB_HOST=localhost
DB_PORT=3306
DB_NAME=portaldata
DB_USER=root
DB_PASSWORD=your_password

# Redis для rate limiting
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 📋 Чек-лист развертывания

- [ ] Применена интеграция тарифов
- [ ] Проверены созданные функции
- [ ] Проверены представления
- [ ] Перезапущен сервис portaldata-api
- [ ] Проверена работоспособность API
- [ ] Протестированы лимиты для разных пользователей
- [ ] Проверены логи на ошибки

## 🎉 Результат

После интеграции:

- ✅ **Rate limiting использует лимиты из тарифов** вместо фиксированных значений
- ✅ **Пользователи связаны с тарифами** через поле `tariff_id`
- ✅ **Лимиты динамические** и зависят от тарифа пользователя
- ✅ **Система отказоустойчива** - при ошибках используются дефолтные лимиты
- ✅ **Полная совместимость** с существующей структурой тарифов
- ✅ **Удобные инструменты** для управления и мониторинга

**Система готова к использованию!** 🚀

---

**Версия**: 1.0  
**Дата**: 2024  
**Статус**: ✅ Интеграция завершена
