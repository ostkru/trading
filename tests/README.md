# 🧪 Тесты Trading API

Эта папка содержит все тесты для Trading API системы.

## 📁 Структура папок

### `api/` - API тесты
- **`test_redis_rate_limiting.php`** - Тесты Redis rate limiting
- **`run_redis_tests.php`** - Запуск Redis тестов

### `integration/` - Интеграционные тесты
- **`comprehensive_improved.php`** - Комплексные интеграционные тесты всех API endpoints

### `unit/` - Модульные тесты
- Модульные тесты для отдельных компонентов

### `performance/` - Тесты производительности
- Тесты нагрузки и производительности

## 🚀 Запуск тестов

### Комплексные тесты
```bash
cd tests/integration
php comprehensive_improved.php
```

### Redis Rate Limiting тесты
```bash
cd tests/api
php test_redis_rate_limiting.php
# или
php run_redis_tests.php
```

### Все тесты
```bash
cd tests
./run_tests.sh
```

## 📊 Результаты тестов

Результаты тестов сохраняются в файлы с расширением `.log` в соответствующих папках.

## 🔧 Требования

- PHP 7.4+
- Redis сервер
- Go API сервис на порту 8095
- MySQL база данных

## 📝 Добавление новых тестов

1. Создайте тест в соответствующей папке по типу
2. Добавьте описание в этот README
3. Обновите `run_tests.sh` если необходимо



