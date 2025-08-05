# Go-Mod API

REST API для управления продуктами, офферами, заказами и складами.

## Описание

Этот проект представляет собой полнофункциональное API, построенное на Go с использованием Gin framework. API предоставляет возможности для:

- Управления продуктами (создание, чтение, обновление, удаление)
- Управления офферами (предложениями)
- Управления заказами
- Управления складами
- Статистики и аналитики

## Технологии

- **Go** - основной язык программирования
- **Gin** - веб-фреймворк
- **MySQL** - база данных
- **OpenAPI/Swagger** - документация API
- **Redoc** - интерактивная документация
- **Nginx** - reverse proxy и SSL termination

## Структура проекта

```
go-mod/
├── cmd/api/           # Точка входа приложения
├── internal/          # Внутренние модули
│   ├── modules/       # Бизнес-логика
│   │   ├── product/   # Модуль продуктов
│   │   ├── offer/     # Модуль офферов
│   │   ├── order/     # Модуль заказов
│   │   ├── warehouse/ # Модуль складов
│   │   └── user/      # Модуль пользователей
│   └── pkg/           # Общие пакеты
│       └── middleware/ # Middleware
├── openapi.json       # OpenAPI спецификация
├── redoc-documentation.html # Документация
├── nginx-trading.conf # Конфигурация nginx
└── comprehensive_api_test.php # Тесты
```

## Установка и запуск

### Требования
- Go 1.19+
- MySQL 8.0+
- Nginx

### 1. Клонирование репозитория
```bash
git clone https://github.com/ostkru/go-mod.git
cd go-mod
```

### 2. Настройка базы данных
```sql
CREATE DATABASE portaldata;
USE portaldata;

-- Таблицы создаются автоматически при первом запуске
```

### 3. Сборка приложения
```bash
go build -o app cmd/api/main.go
```

### 4. Запуск в режиме разработки
```bash
./app
```

Приложение будет доступно по адресу: `http://localhost:8095`

### 5. Настройка nginx для продакшена

1. Скопируйте конфигурацию:
```bash
sudo cp nginx-trading.conf /etc/nginx/sites-available/portaldata-trading
```

2. Активируйте сайт:
```bash
sudo ln -s /etc/nginx/sites-available/portaldata-trading /etc/nginx/sites-enabled/
```

3. Проверьте конфигурацию:
```bash
sudo nginx -t
```

4. Перезапустите nginx:
```bash
sudo systemctl reload nginx
```

## API Документация

### Продакшен
- **Базовый URL**: `https://api.portaldata.ru/v1/trading`
- **Документация**: `https://api.portaldata.ru/v1/trading/docs`

### Разработка
- **Базовый URL**: `http://localhost:8095/api/v1`
- **Документация**: `http://localhost:8095/redoc-documentation.html`

## Аутентификация

Все API методы требуют аутентификации через API ключ:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     https://api.portaldata.ru/v1/trading/products
```

## Основные эндпоинты

### Продукты
- `GET /products` - Список продуктов
- `POST /products` - Создание продукта
- `GET /products/{id}` - Получение продукта
- `PUT /products/{id}` - Обновление продукта
- `DELETE /products/{id}` - Удаление продукта

### Офферы
- `GET /offers` - Список офферов
- `POST /offers` - Создание оффера
- `GET /offers/{id}` - Получение оффера
- `PUT /offers/{id}` - Обновление оффера
- `DELETE /offers/{id}` - Удаление оффера
- `GET /offers/public` - Публичные офферы

### Заказы
- `GET /orders` - Список заказов
- `POST /orders` - Создание заказа
- `GET /orders/{id}` - Получение заказа
- `PUT /orders/{id}` - Обновление заказа
- `DELETE /orders/{id}` - Удаление заказа

### Склады
- `GET /warehouses` - Список складов
- `POST /warehouses` - Создание склада
- `GET /warehouses/{id}` - Получение склада
- `PUT /warehouses/{id}` - Обновление склада
- `DELETE /warehouses/{id}` - Удаление склада

### Статистика
- `GET /statistics` - Статистика пользователя

## Тестирование

### Запуск комплексного теста
```bash
php comprehensive_api_test.php
```

### Тестирование продакшена
```bash
# Обновите URL в тестах на продакшен
sed -i 's|http://localhost:8095/api/v1|https://api.portaldata.ru/v1/trading|g' comprehensive_api_test.php
php comprehensive_api_test.php
```

## Архитектура

### Reverse Proxy (Nginx)
- SSL termination
- CORS поддержка
- Логирование
- Кэширование (опционально)

### Go Application
- REST API на Gin
- Middleware для аутентификации
- Middleware для статистики
- Валидация данных

### База данных
- MySQL для хранения данных
- Автоматическое создание таблиц
- Индексы для оптимизации

## Мониторинг

### Логи nginx
```bash
tail -f /var/log/nginx/portaldata-trading-access.log
tail -f /var/log/nginx/portaldata-trading-error.log
```

### Логи приложения
```bash
tail -f app.log
```

### Статус сервисов
```bash
sudo systemctl status nginx
ps aux | grep app
```

## Развертывание

### Автоматический деплой
```bash
# Остановка старого процесса
pkill -f "./app"

# Сборка нового бинарного файла
go build -o app cmd/api/main.go

# Запуск нового процесса
./app > app.log 2>&1 &

# Проверка статуса
netstat -tlnp | grep 8090
```

### Обновление документации
```bash
redoc-cli bundle openapi.json -o redoc-documentation.html
sudo systemctl reload nginx
```

## Безопасность

- Все запросы требуют аутентификации
- SSL/TLS шифрование
- Валидация входных данных
- Проверка прав доступа к ресурсам
- Логирование всех операций

## Поддержка

Для вопросов и предложений создавайте issues в репозитории:
https://github.com/ostkru/go-mod/issues 