package main

import (
	"log"
	"time"

	"portaldata-api/internal/modules/ratelimit"

	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем Redis-based rate limit сервис
	redisService := ratelimit.NewRedisRateLimitService(
		"127.0.0.1:6379", // адрес Redis
		"",               // пароль (пустой)
		0,                // база данных
	)

	if redisService == nil {
		log.Fatal("Не удалось подключиться к Redis")
	}
	defer redisService.Close()

	// Создаем handler для управления API ключами
	redisHandler := ratelimit.NewRedisHandler(redisService)

	// Создаем Gin роутер
	r := gin.Default()

	// Настраиваем middleware для rate limiting
	r.Use(ratelimit.RedisRateLimitMiddleware(redisService))

	// Группа для rate limit API
	rateLimitAPI := r.Group("/api/v1/rate-limit")
	ratelimit.SetupRedisRoutes(rateLimitAPI, redisHandler)

	// Примеры эндпоинтов с rate limiting
	offers := r.Group("/api/v1/offers")
	{
		// Публичные офферы с отдельным лимитом
		offers.GET("/public", func(c *gin.Context) {
			apiKey, _ := c.Get("apiKey")
			c.JSON(200, gin.H{
				"message":   "Публичные офферы",
				"api_key":   apiKey,
				"timestamp": time.Now(),
			})
		})

		// Обычные офферы с общим лимитом
		offers.GET("", func(c *gin.Context) {
			apiKey, _ := c.Get("apiKey")
			c.JSON(200, gin.H{
				"message":   "Список офферов",
				"api_key":   apiKey,
				"timestamp": time.Now(),
			})
		})
	}

	// Запускаем сервер
	log.Println("Сервер запущен на :8080")
	log.Println("Rate limiting API доступен на /api/v1/rate-limit")
	log.Println("Примеры запросов:")
	log.Println("  GET /api/v1/rate-limit/stats - общая статистика")
	log.Println("  GET /api/v1/rate-limit/search?pattern=* - поиск API ключей")
	log.Println("  GET /api/v1/rate-limit/api-keys/YOUR_API_KEY - информация об API ключе")
	log.Println("  GET /api/v1/rate-limit/top?limit=5 - топ API ключей")

	r.Run(":8080")
}

// Пример использования API для мониторинга:
/*
# Получить общую статистику
curl http://localhost:8080/api/v1/rate-limit/stats

# Поиск API ключей
curl http://localhost:8080/api/v1/rate-limit/search?pattern=user*

# Информация об API ключе
curl http://localhost:8080/api/v1/rate-limit/api-keys/YOUR_API_KEY

# Статистика по API ключу
curl http://localhost:8080/api/v1/rate-limit/api-keys/YOUR_API_KEY/stats

# Топ API ключей
curl http://localhost:8080/api/v1/rate-limit/top?limit=10

# Сброс счетчиков
curl -X POST http://localhost:8080/api/v1/rate-limit/api-keys/YOUR_API_KEY/reset?type=all

# Тест rate limiting
curl -H "X-API-KEY: YOUR_API_KEY" http://localhost:8080/api/v1/offers/public
*/
