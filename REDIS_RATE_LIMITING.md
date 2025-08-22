# Redis-based Rate Limiting для Trading API

## Обзор

Система rate limiting на основе Redis для управления лимитами API запросов с возможностью поиска по API ключам и мониторинга использования.

## Особенности

- **Минутные лимиты**: 60 запросов в минуту для всех методов
- **Дневные лимиты**: 1000 запросов в день для GET методов
- **Раздельные лимиты**: Отдельные лимиты для публичных и приватных эндпоинтов
- **Поиск по API ключам**: Возможность поиска и мониторинга API ключей
- **Автоматический сброс**: Счетчики автоматически сбрасываются по времени
- **Redis-based**: Высокая производительность и масштабируемость

## Архитектура

### Компоненты

1. **RedisRateLimitService** - основной сервис для работы с Redis
2. **RedisRateLimitMiddleware** - middleware для проверки лимитов
3. **RedisHandler** - API endpoints для управления и мониторинга
4. **Конфигурация** - настройки Redis и лимитов

### Структура ключей Redis

```
rate_limit:{API_KEY}:{ENDPOINT_TYPE}:{COUNTER_TYPE}
```

Примеры:
- `rate_limit:user123:all:minute` - минутный счетчик для всех эндпоинтов
- `rate_limit:user123:public:day` - дневной счетчик для публичных эндпоинтов
- `rate_limit:user123:all:last_request` - время последнего запроса

## Установка и настройка

### 1. Зависимости

```bash
go get github.com/go-redis/redis/v8
```

### 2. Конфигурация Redis

```yaml
# configs/redis.yaml
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
  pool_size: 10

rate_limits:
  minute:
    default: 60
    public: 120
  day:
    default: 1000
    public: 2000
```

### 3. Инициализация

```go
import "portaldata-api/internal/modules/ratelimit"

// Создаем Redis сервис
redisService := ratelimit.NewRedisRateLimitService(
    "127.0.0.1:6379",
    "",
    0,
)

// Создаем middleware
r.Use(ratelimit.RedisRateLimitMiddleware(redisService))

// Настраиваем API endpoints
redisHandler := ratelimit.NewRedisHandler(redisService)
ratelimit.SetupRedisRoutes(r.Group("/api/v1/rate-limit"), redisHandler)
```

## API Endpoints

### Мониторинг и управление

#### 1. Общая статистика
```http
GET /api/v1/rate-limit/stats
```

Ответ:
```json
{
  "total_api_keys": 150,
  "total_keys": 450,
  "redis_info": "...",
  "active_api_keys": ["user1", "user2", ...]
}
```

#### 2. Поиск API ключей
```http
GET /api/v1/rate-limit/search?pattern=user*
```

Ответ:
```json
{
  "pattern": "user*",
  "api_keys": ["user1", "user2", "user123"],
  "count": 3
}
```

#### 3. Информация об API ключе
```http
GET /api/v1/rate-limit/api-keys/{API_KEY}
```

Ответ:
```json
{
  "api_key": "user123",
  "endpoints": {
    "all": {
      "minute_count": 15,
      "day_count": 150,
      "last_request": "2024-01-15T10:30:00Z"
    },
    "public": {
      "minute_count": 8,
      "day_count": 80,
      "last_request": "2024-01-15T10:25:00Z"
    }
  },
  "total_requests": 230,
  "last_request": "2024-01-15T10:30:00Z"
}
```

#### 4. Статистика по API ключу
```http
GET /api/v1/rate-limit/api-keys/{API_KEY}/stats
```

#### 5. Топ API ключей
```http
GET /api/v1/rate-limit/top?limit=10
```

#### 6. Сброс счетчиков
```http
POST /api/v1/rate-limit/api-keys/{API_KEY}/reset?type=all
```

Параметры `type`:
- `all` - сброс всех счетчиков
- `public` - сброс только публичных счетчиков
- `all` - сброс только общих счетчиков

## Использование в приложении

### 1. Middleware

```go
// Применяем ко всем маршрутам
r.Use(ratelimit.RedisRateLimitMiddleware(redisService))

// Или к конкретной группе
offers := r.Group("/offers")
offers.Use(ratelimit.RedisRateLimitMiddleware(redisService))
```

### 2. Получение API ключа в handler

```go
func GetOffers(c *gin.Context) {
    apiKey, exists := c.Get("apiKey")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "API ключ не найден"})
        return
    }
    
    // Используем apiKey для логирования или других целей
    log.Printf("Запрос от API ключа: %s", apiKey)
    
    // ... логика получения офферов
}
```

### 3. Заголовки ответа

Система автоматически добавляет заголовки с информацией о лимитах:

```
X-RateLimit-Limit-Minute: 60
X-RateLimit-Limit-Day: 1000
X-RateLimit-Remaining-Minute: 45
X-RateLimit-Remaining-Day: 850
```

## Мониторинг и отладка

### 1. Redis CLI

```bash
# Подключение к Redis
redis-cli

# Просмотр всех ключей rate limiting
KEYS "rate_limit:*"

# Получение значения счетчика
GET "rate_limit:user123:all:minute"

# Просмотр TTL ключа
TTL "rate_limit:user123:all:minute"

# Удаление ключа (сброс счетчика)
DEL "rate_limit:user123:all:minute"
```

### 2. Логирование

```go
// Включение детального логирования
log.SetLevel(log.DebugLevel)

// Логирование в middleware
log.Printf("API ключ: %s, эндпоинт: %s, разрешен: %v", 
    apiKey, endpoint, rateLimitCheck.Allowed)
```

### 3. Метрики

```go
// Получение статистики Redis
stats := redisService.GetRateLimitStats()
log.Printf("Активных API ключей: %d", stats.TotalAPIKeys)
log.Printf("Всего ключей в Redis: %d", stats.TotalKeys)
```

## Производительность

### 1. Redis операции

- **INCR** - атомарное увеличение счетчика
- **EXPIRE** - автоматическое удаление по TTL
- **PIPELINE** - группировка операций для производительности
- **SCAN** - безопасный поиск ключей

### 2. Оптимизации

- Ключи автоматически удаляются по TTL
- Использование pipeline для атомарных операций
- Кэширование статистики в памяти
- Асинхронное обновление счетчиков

## Безопасность

### 1. Валидация API ключей

- Проверка формата ключа
- Логирование всех запросов
- Ограничение длины ключа

### 2. Rate limiting bypass

- Административные ключи с повышенными лимитами
- Whitelist для определенных IP адресов
- Emergency override для критических операций

## Миграция с MySQL

### 1. Параллельное использование

```go
// Можно использовать оба сервиса одновременно
mysqlService := ratelimit.NewService(db)
redisService := ratelimit.NewRedisRateLimitService(redisAddr, "", 0)

// Выбор сервиса по конфигурации
if useRedis {
    middleware = ratelimit.RedisRateLimitMiddleware(redisService)
} else {
    middleware = ratelimit.RateLimitMiddleware(mysqlService)
}
```

### 2. Синхронизация данных

```go
// Экспорт данных из MySQL в Redis
func SyncToRedis(mysqlService *Service, redisService *RedisRateLimitService) error {
    // Получение всех записей из MySQL
    // Создание соответствующих ключей в Redis
    // Установка TTL для автоматического удаления
}
```

## Troubleshooting

### 1. Redis недоступен

```go
// Fallback на MySQL
if redisService == nil {
    log.Println("Redis недоступен, используем MySQL")
    middleware = ratelimit.RateLimitMiddleware(mysqlService)
}
```

### 2. Высокое потребление памяти

```bash
# Проверка размера ключей
redis-cli --bigkeys

# Очистка старых ключей
redis-cli FLUSHDB
```

### 3. Медленные запросы

```go
// Увеличение размера pipeline
pipe := s.client.Pipeline()
pipe.Incr(s.ctx, key)
pipe.Expire(s.ctx, key, ttl)
_, err := pipe.Exec(s.ctx)
```

## Примеры использования

### 1. Базовый rate limiting

```go
func main() {
    r := gin.Default()
    
    // Создаем Redis сервис
    redisService := ratelimit.NewRedisRateLimitService("localhost:6379", "", 0)
    
    // Применяем middleware
    r.Use(ratelimit.RedisRateLimitMiddleware(redisService))
    
    // API endpoints
    r.GET("/api/offers", GetOffers)
    r.GET("/api/offers/public", GetPublicOffers)
    
    r.Run(":8080")
}
```

### 2. Мониторинг в реальном времени

```go
func MonitorRateLimits(redisService *RedisRateLimitService) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        stats, err := redisService.GetRateLimitStats()
        if err != nil {
            log.Printf("Ошибка получения статистики: %v", err)
            continue
        }
        
        log.Printf("Активных API ключей: %d, общих ключей: %d", 
            stats.TotalAPIKeys, stats.TotalKeys)
    }
}
```

### 3. Автоматический сброс счетчиков

```go
func AutoResetCounters(redisService *RedisRateLimitService) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        // Получаем все API ключи
        apiKeys, err := redisService.SearchAPIKeys("*")
        if err != nil {
            continue
        }
        
        // Сбрасываем дневные счетчики
        for _, apiKey := range apiKeys {
            redisService.ResetCounters(apiKey, "all")
        }
    }
}
```

## Заключение

Redis-based rate limiting система предоставляет:

- Высокую производительность
- Масштабируемость
- Гибкость настройки
- Мониторинг в реальном времени
- Автоматическое управление ресурсами

Система готова к продакшену и может быть легко интегрирована в существующее приложение.
