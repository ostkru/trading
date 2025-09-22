package ratelimit

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RedisRateLimitMiddleware создает middleware для проверки лимитов через Redis
func RedisRateLimitMiddleware(redisService *RedisRateLimitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Полностью публичные маршруты не требуют rate limiting
		if c.Request.URL.Path == "/" ||
			c.Request.URL.Path == "/browser" ||
			strings.Contains(c.Request.URL.Path, "/swagger") ||
			strings.Contains(c.Request.URL.Path, "/rate-limit") {
			c.Next()
			return
		}

		// Объявляем переменную для API ключа
		var apiKey string

		// Публичные офферы не требуют API ключа, но проходят rate limiting
		if strings.Contains(c.Request.URL.Path, "/offers/public") {
			// Для публичных офферов используем специальный API ключ для rate limiting
			apiKey = "public_offers"
			// Продолжаем выполнение для применения rate limiting
		}

		// Получаем API ключ из различных источников (если еще не установлен для публичных офферов)
		if apiKey == "" {
			// Проверяем заголовок Authorization
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = authHeader[7:]
			} else {
				// Проверяем заголовок X-API-KEY
				apiKey = c.GetHeader("X-API-KEY")
			}

			// Если токен не найден в заголовках, проверяем GET параметр api_key
			if apiKey == "" {
				apiKey = c.Query("api_key")
			}
		}

		// Специальная обработка для запросов без API ключа
		// Rate limiting должен работать для всех запросов, включая неавторизованные
		if apiKey == "" {
			// Проверяем, является ли это тестовым запросом
			userAgent := c.GetHeader("User-Agent")
			if strings.Contains(userAgent, "PHP") ||
				strings.Contains(userAgent, "curl") ||
				strings.Contains(userAgent, "Postman") ||
				c.GetHeader("X-Test-Request") == "true" {
				// Для тестовых запросов используем специальный API ключ
				apiKey = "test_request"
			} else {
				// Для обычных запросов без API ключа используем IP адрес
				apiKey = "anonymous_" + c.ClientIP()
			}
		}

		// Определяем тип эндпоинта для лимитов
		endpoint := "all" // по умолчанию для всех методов

		// Для публичных методов используем отдельный лимит
		// Проверяем, является ли это публичным endpoint (без авторизации)
		if strings.Contains(c.Request.URL.Path, "/offers/public") ||
			strings.Contains(c.Request.URL.Path, "/api/products") {
			endpoint = "public"
		}

		// Определяем, является ли запрос GET методом
		isGetMethod := c.Request.Method == "GET"

		// Проверяем лимиты
		rateLimitCheck, err := redisService.CheckRateLimit(apiKey, endpoint, isGetMethod)
		if err != nil {
			log.Printf("Ошибка проверки лимитов: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки лимитов"})
			c.Abort()
			return
		}

		// Если лимит превышен, возвращаем ошибку
		if !rateLimitCheck.Allowed {
			c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", rateLimitCheck.MinuteLimit))
			c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", rateLimitCheck.DayLimit))
			c.Header("X-RateLimit-Remaining-Minute", fmt.Sprintf("%d", rateLimitCheck.MinuteLimit-rateLimitCheck.MinuteUsed))
			c.Header("X-RateLimit-Remaining-Day", fmt.Sprintf("%d", rateLimitCheck.DayLimit-rateLimitCheck.DayUsed))
			c.Header("X-RateLimit-Reset-Minute", fmt.Sprintf("%d", 60))    // секунды до сброса
			c.Header("X-RateLimit-Reset-Day", fmt.Sprintf("%d", 24*60*60)) // секунды до сброса

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Превышен лимит запросов",
				"message": rateLimitCheck.Message,
				"limits": gin.H{
					"minute": gin.H{
						"limit":     rateLimitCheck.MinuteLimit,
						"used":      rateLimitCheck.MinuteUsed,
						"remaining": rateLimitCheck.MinuteLimit - rateLimitCheck.MinuteUsed,
					},
					"day": gin.H{
						"limit":     rateLimitCheck.DayLimit,
						"used":      rateLimitCheck.DayUsed,
						"remaining": rateLimitCheck.DayLimit - rateLimitCheck.DayUsed,
					},
				},
			})
			c.Abort()
			return
		}

		// Добавляем информацию о лимитах в заголовки для всех запросов
		c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", rateLimitCheck.MinuteLimit))
		c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", rateLimitCheck.DayLimit))
		c.Header("X-RateLimit-Remaining-Minute", fmt.Sprintf("%d", rateLimitCheck.MinuteLimit-rateLimitCheck.MinuteUsed))
		c.Header("X-RateLimit-Remaining-Day", fmt.Sprintf("%d", rateLimitCheck.DayLimit-rateLimitCheck.DayUsed))
		c.Header("X-RateLimit-Reset-Minute", fmt.Sprintf("%d", 60))    // секунды до сброса
		c.Header("X-RateLimit-Reset-Day", fmt.Sprintf("%d", 24*60*60)) // секунды до сброса

		// Добавляем API ключ в контекст для дальнейшего использования
		c.Set("apiKey", apiKey)

		c.Next()
	}
}
