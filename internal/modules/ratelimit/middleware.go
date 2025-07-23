package ratelimit

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware создает middleware для проверки лимитов
func RateLimitMiddleware(rateLimitService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем API ключ из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API ключ не предоставлен"})
			c.Abort()
			return
		}

		// Извлекаем токен из Bearer
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат API ключа"})
			c.Abort()
			return
		}

		// Получаем userID из контекста (должен быть установлен в auth middleware)
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ID пользователя"})
			c.Abort()
			return
		}

		// Определяем тип эндпоинта для лимитов
		endpoint := "all" // по умолчанию для всех методов
		
		// Для публичных методов используем отдельный лимит
		if strings.Contains(c.Request.URL.Path, "/offers/public") {
			endpoint = "public"
		}

		// Определяем, является ли запрос GET методом
		isGetMethod := c.Request.Method == "GET"
		
		// Добавляем логирование для отладки
		log.Printf("RateLimit: userID=%d, endpoint=%s, path=%s, method=%s, isGet=%v", 
			userID, endpoint, c.Request.URL.Path, c.Request.Method, isGetMethod)

		// Проверяем лимиты
		check, err := rateLimitService.CheckRateLimit(userID, token, endpoint, isGetMethod)
		if err != nil {
			log.Printf("RateLimit error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки лимитов"})
			c.Abort()
			return
		}

		log.Printf("RateLimit check: allowed=%v, minute_used=%d/%d, day_used=%d/%d", 
			check.Allowed, check.MinuteUsed, check.MinuteLimit, check.DayUsed, check.DayLimit)

		// Если лимит превышен, возвращаем 429
		if !check.Allowed {
			c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", check.MinuteLimit))
			c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", check.DayLimit))
			c.Header("X-RateLimit-Remaining-Minute", fmt.Sprintf("%d", check.MinuteLimit-check.MinuteUsed))
			c.Header("X-RateLimit-Remaining-Day", fmt.Sprintf("%d", check.DayLimit-check.DayUsed))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": check.Message,
				"limits": gin.H{
					"minute_limit": check.MinuteLimit,
					"day_limit":    check.DayLimit,
					"minute_used":  check.MinuteUsed,
					"day_used":     check.DayUsed,
				},
			})
			c.Abort()
			return
		}

		// Добавляем заголовки с информацией о лимитах
		c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", check.MinuteLimit))
		c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", check.DayLimit))
		c.Header("X-RateLimit-Remaining-Minute", fmt.Sprintf("%d", check.MinuteLimit-check.MinuteUsed))
		c.Header("X-RateLimit-Remaining-Day", fmt.Sprintf("%d", check.DayLimit-check.DayUsed))

		c.Next()
	}
} 