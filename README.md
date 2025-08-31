# 🚀 Trading API - Торговая платформа

Полнофункциональная API для торговой платформы с автоматической генерацией документации.

## ✨ Возможности

- **Продукты** - управление каталогом товаров
- **Предложения** - торговые предложения покупки/продажи
- **Склады** - управление складскими запасами
- **Заказы** - обработка заказов клиентов
- **Rate Limiting** - защита от перегрузки API
- **Автоматическая документация** - генерация OpenAPI спецификации

## 🏗️ Архитектура

```
├── cmd/
│   ├── api/                    # Основное API приложение
│   └── openapi-generator/      # Генератор OpenAPI документации
├── internal/
│   ├── modules/                # Бизнес-логика
│   │   ├── products/          # Управление продуктами
│   │   ├── offers/            # Управление предложениями
│   │   ├── warehouses/        # Управление складами
│   │   ├── orders/            # Управление заказами
│   │   └── ratelimit/         # Rate limiting
│   └── pkg/                   # Общие пакеты
├── configs/                    # Конфигурационные файлы
├── scripts/                    # Скрипты управления
└── docs/                       # Документация
```

## 🚀 Быстрый старт

### 1. Запуск сервисов

```bash
# Установка и запуск всех сервисов
make install-services

# Или по отдельности
make start-redis
make start-api
```

### 2. Генерация документации

```bash
# Автоматическая генерация OpenAPI документации
make generate-docs

# Или вручную
make openapi-generator
./openapi-generator
```

### 3. Тестирование

```bash
# Комплексное тестирование
php comprehensive_improved.php

# Тестирование Redis Rate Limiting
php test_redis_rate_limiting.php
```

## 📚 Автоматический генератор OpenAPI документации

### 🎯 Возможности

- **Автоматический анализ Go кода** - сканирует структуры и обработчики
- **Генерация OpenAPI 3.0.3** - полная совместимость со стандартом
- **Автоматические примеры** - создает примеры запросов и ответов
- **HTML документация** - генерирует красивую веб-страницу
- **Валидация тегов** - анализирует binding и json теги
- **Безопасность** - автоматически добавляет схемы аутентификации

### 🔍 Как это работает

#### 1. Сканирование Go файлов

Генератор автоматически сканирует директорию `internal/modules/` и анализирует:

- **Структуры** - для генерации схем данных
- **Обработчики** - для генерации API endpoints
- **Комментарии** - для описаний и метаданных
- **Теги** - для валидации и требований

#### 2. Анализ AST

```go
// Пример комментария для автоматического определения HTTP метода и пути
// @POST /products
func (h *Handler) CreateProduct(c *gin.Context) {
    // ...
}
```

#### 3. Генерация схем

```go
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required"`
    Article     string  `json:"article" binding:"required"`
    Brand       string  `json:"brand" binding:"required"`
    Category    string  `json:"category" binding:"required"`
    Price       float64 `json:"price" binding:"required,min=0"`
    Description string  `json:"description"`
}
```

Автоматически генерируется:

```json
{
  "CreateProductRequest": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "description": "Название продукта"
      },
      "price": {
        "type": "number",
        "format": "double",
        "minimum": 0,
        "example": 123.45
      }
    },
    "required": ["name", "article", "brand", "category", "price"]
  }
}
```

#### 4. Генерация примеров

Автоматически создаются примеры на основе:

- **Типов данных** - string, int, float64, bool
- **Названий полей** - name, email, url
- **Валидации** - min, max, email, url
- **Структуры** - вложенные объекты

### 📊 Выходные файлы

#### 1. `openapi_generated.json`

Полная OpenAPI 3.0.3 спецификация с примерами запросов и ответов.

#### 2. `api_documentation.html`

Красивая HTML документация с:
- **Цветовая кодировка** HTTP методов
- **Примеры запросов** и ответов
- **Схемы данных** с описаниями
- **Адаптивный дизайн** для мобильных устройств

### 🛠️ Использование

```bash
# Сборка генератора
make openapi-generator

# Генерация документации
make generate-docs

# Очистка
make clean-docs
```

## 🔧 Управление сервисами

### Redis

```bash
# Запуск
make start-redis

# Остановка
make stop-redis

# Статус
make status-redis

# Перезапуск
make restart-redis
```

### Trading API

```bash
# Запуск
make start-api

# Остановка
make stop-api

# Статус
make status-api

# Перезапуск
make restart-api
```

### Все сервисы

```bash
# Установка
make install-services

# Запуск всех
make start-all

# Остановка всех
make stop-all

# Перезапуск всех
make restart-all

# Удаление
make uninstall-services
```

## 📡 API Endpoints

### Продукты

- `POST /products` - создание продукта
- `GET /products` - список продуктов
- `GET /products/{id}` - получение продукта
- `PUT /products/{id}` - обновление продукта
- `DELETE /products/{id}` - удаление продукта

### Предложения

- `POST /offers` - создание предложения
- `GET /offers` - список предложений
- `GET /offers/public` - публичные предложения
- `PUT /offers/{id}` - обновление предложения
- `DELETE /offers/{id}` - удаление предложения

### Склады

- `POST /warehouses` - создание склада
- `GET /warehouses` - список складов
- `POST /warehouses/batch` - пакетное создание складов
- `PUT /warehouses/{id}` - обновление склада
- `DELETE /warehouses/{id}` - удаление склада

### Заказы

- `POST /orders` - создание заказа
- `GET /orders` - список заказов
- `GET /orders/{id}` - получение заказа
- `PUT /orders/{id}` - обновление заказа
- `DELETE /orders/{id}` - удаление заказа

### Rate Limiting

- `GET /rate-limit/stats` - статистика использования
- `GET /rate-limit/info` - информация об API ключе
- `POST /rate-limit/reset` - сброс счетчиков
- `GET /rate-limit/top` - топ API ключей

## 🔐 Аутентификация

API использует API ключи для аутентификации:

```bash
curl -H "X-API-KEY: your_api_key" \
     https://api.portaldata.ru/v1/trading/products
```

## 📊 Rate Limiting

- **Публичные endpoints**: 2000 запросов/минуту, 20000/день
- **Защищенные endpoints**: 1000 запросов/минуту, 10000/день
- **Лимиты настраиваются** через `configs/redis.yaml`

## 🧪 Тестирование

### Комплексное тестирование

```bash
php comprehensive_improved.php
```

### Тестирование Redis Rate Limiting

```bash
php test_redis_rate_limiting.php
```

### Запуск отдельных тестов

```bash
# Тестирование продуктов
php -f comprehensive_improved.php -- --test=products

# Тестирование предложений
php -f comprehensive_improved.php -- --test=offers
```

## 📁 Структура проекта

```
├── cmd/                        # Исполняемые файлы
├── internal/                   # Внутренняя логика
├── configs/                    # Конфигурация
├── scripts/                    # Скрипты управления
├── docs/                       # Документация
├── examples/                   # Примеры использования
├── Makefile                    # Команды управления
├── README.md                   # Основная документация
├── OPENAPI_GENERATOR_README.md # Документация генератора
├── SYSTEMD_SERVICES_README.md  # Документация сервисов
└── QUICK_START.md             # Быстрый старт
```

## 🚀 Развертывание

### Локальная разработка

```bash
# Клонирование
git clone <repository>
cd trading

# Установка зависимостей
go mod download

# Запуск сервисов
make install-services
make start-all

# Генерация документации
make generate-docs
```

### Продакшн

```bash
# Сборка
make build

# Установка сервисов
make install-services

# Запуск
make start-all

# Проверка статуса
make status-all
```

## 🔧 Конфигурация

### Redis

```yaml
# configs/redis.yaml
redis:
  host: localhost
  port: 6379
  db: 0

rate_limits:
  minute:
    default: 1000
    public: 2000
  day:
    default: 10000
    public: 20000
```

### Nginx

```nginx
# /etc/nginx/sites-available/api.portaldata.ru
location /v1/trading/ {
    proxy_pass http://localhost:8095/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

## 📚 Документация

- [OpenAPI Generator README](OPENAPI_GENERATOR_README.md) - Подробная документация генератора
- [Systemd Services README](SYSTEMD_SERVICES_README.md) - Управление сервисами
- [Quick Start Guide](QUICK_START.md) - Быстрый старт

## 🤝 Вклад в проект

1. **Fork** репозитория
2. **Создайте** feature branch (`git checkout -b feature/amazing-feature`)
3. **Сделайте** коммит (`git commit -m 'Add amazing feature'`)
4. **Push** в branch (`git push origin feature/amazing-feature`)
5. **Создайте** Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 📞 Поддержка

Если у вас есть вопросы или предложения:

1. **Создайте Issue** в репозитории
2. **Напишите в чат** команды разработки
3. **Проверьте документацию** по ссылкам выше

---

**🎉 Автоматизируйте документацию и сосредоточьтесь на коде!** 