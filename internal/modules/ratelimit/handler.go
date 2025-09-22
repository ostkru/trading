package ratelimit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RedisHandler struct {
	redisService *RedisRateLimitService
}

func NewRedisHandler(redisService *RedisRateLimitService) *RedisHandler {
	return &RedisHandler{
		redisService: redisService,
	}
}

// GetAPIKeyInfo получает информацию об API ключе
func (h *RedisHandler) GetAPIKeyInfo(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ не указан"})
		return
	}

	info, err := h.redisService.GetAPIKeyInfo(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения информации об API ключе"})
		return
	}

	c.JSON(http.StatusOK, info)
}

// SearchAPIKeys поиск API ключей по паттерну
func (h *RedisHandler) SearchAPIKeys(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		pattern = "*" // По умолчанию ищем все ключи
	}

	apiKeys, err := h.redisService.SearchAPIKeys(pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска API ключей"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pattern":  pattern,
		"api_keys": apiKeys,
		"count":    len(apiKeys),
	})
}

// GetRateLimitStats получает общую статистику rate limiting
func (h *RedisHandler) GetRateLimitStats(c *gin.Context) {
	stats, err := h.redisService.GetRateLimitStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения статистики"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAPIKeyStats получает статистику по конкретному API ключу
func (h *RedisHandler) GetAPIKeyStats(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ не указан"})
		return
	}

	// Получаем детальную статистику
	info, err := h.redisService.GetAPIKeyInfo(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения статистики API ключа"})
		return
	}

	// Получаем лимиты из тарифа пользователя
	minuteLimit, dayLimit, err := h.redisService.getUserLimitsFromAPIKey(apiKey)
	if err != nil {
		// Если не удалось получить лимиты, используем дефолтные
		minuteLimit = 60
		dayLimit = 1000
	}

	// Формируем расширенную статистику
	stats := gin.H{
		"api_key":        apiKey,
		"total_requests": info.TotalRequests,
		"last_request":   info.LastRequest,
		"endpoints":      info.Endpoints,
		"limits": gin.H{
			"minute": gin.H{
				"limit":       minuteLimit,
				"description": "Запросов в минуту",
			},
			"day": gin.H{
				"limit":       dayLimit,
				"description": "Запросов в день (только GET)",
			},
		},
	}

	c.JSON(http.StatusOK, stats)
}

// ResetAPIKeyCounters сбрасывает счетчики для API ключа
func (h *RedisHandler) ResetAPIKeyCounters(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ не указан"})
		return
	}

	// Получаем тип сброса
	resetType := c.Query("type")
	if resetType == "" {
		resetType = "all" // По умолчанию сбрасываем все
	}

	// Сбрасываем счетчики
	err := h.redisService.ResetCounters(apiKey, resetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сброса счетчиков"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Счетчики успешно сброшены",
		"api_key":    apiKey,
		"reset_type": resetType,
	})
}

// GetTopAPIKeys получает топ API ключей по количеству запросов
func (h *RedisHandler) GetTopAPIKeys(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Получаем все API ключи
	allKeys, err := h.redisService.SearchAPIKeys("*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка API ключей"})
		return
	}

	// Получаем статистику по каждому ключу
	type APIKeyUsage struct {
		APIKey        string `json:"api_key"`
		TotalRequests int    `json:"total_requests"`
		MinuteCount   int    `json:"minute_count"`
		DayCount      int    `json:"day_count"`
	}

	var usageStats []APIKeyUsage
	for _, apiKey := range allKeys {
		info, err := h.redisService.GetAPIKeyInfo(apiKey)
		if err != nil {
			continue // Пропускаем проблемные ключи
		}

		usageStats = append(usageStats, APIKeyUsage{
			APIKey:        apiKey,
			TotalRequests: info.TotalRequests,
			MinuteCount:   info.Endpoints["all"].MinuteCount,
			DayCount:      info.Endpoints["all"].DayCount,
		})
	}

	// Сортируем по общему количеству запросов (убывание)
	// Простая сортировка пузырьком для небольшого количества ключей
	for i := 0; i < len(usageStats)-1; i++ {
		for j := 0; j < len(usageStats)-i-1; j++ {
			if usageStats[j].TotalRequests < usageStats[j+1].TotalRequests {
				usageStats[j], usageStats[j+1] = usageStats[j+1], usageStats[j]
			}
		}
	}

	// Ограничиваем результат
	if len(usageStats) > limit {
		usageStats = usageStats[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"top_api_keys": usageStats,
		"limit":        limit,
		"total_found":  len(allKeys),
	})
}
