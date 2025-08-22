package ratelimit

import (
	"github.com/gin-gonic/gin"
)

// SetupRedisRoutes настраивает маршруты для Redis-based rate limiting
func SetupRedisRoutes(r *gin.RouterGroup, redisHandler *RedisHandler) {
	// Группа для управления API ключами
	apiKeys := r.Group("/api-keys")
	{
		// Получить информацию об API ключе
		apiKeys.GET("/:api_key", redisHandler.GetAPIKeyInfo)

		// Получить статистику по API ключу
		apiKeys.GET("/:api_key/stats", redisHandler.GetAPIKeyStats)

		// Сбросить счетчики для API ключа
		apiKeys.POST("/:api_key/reset", redisHandler.ResetAPIKeyCounters)
	}

	// Поиск API ключей
	r.GET("/search", redisHandler.SearchAPIKeys)

	// Общая статистика
	r.GET("/stats", redisHandler.GetRateLimitStats)

	// Топ API ключей по использованию
	r.GET("/top", redisHandler.GetTopAPIKeys)
}
