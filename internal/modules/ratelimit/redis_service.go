package ratelimit

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRateLimitService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRateLimitService(redisAddr, redisPassword string, redisDB int) *RedisRateLimitService {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
		PoolSize: 10,
	})

	ctx := context.Background()

	// Проверяем соединение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Ошибка подключения к Redis: %v", err)
		return nil
	}

	return &RedisRateLimitService{
		client: client,
		ctx:    ctx,
	}
}

// CheckRateLimit проверяет лимиты для API ключа
func (s *RedisRateLimitService) CheckRateLimit(apiKey, endpoint string, isGetMethod bool) (*RateLimitCheck, error) {
	// Определяем тип лимита на основе эндпоинта
	limitType := "all"
	if endpoint == "public" {
		limitType = "public"
	}

	// Ключи для Redis
	minuteKey := fmt.Sprintf("rate_limit:%s:%s:minute", apiKey, limitType)
	dayKey := fmt.Sprintf("rate_limit:%s:%s:day", apiKey, limitType)

	// Получаем текущие счетчики
	minuteCount, err := s.getCounter(minuteKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения минутного счетчика: %v", err)
	}

	dayCount, err := s.getCounter(dayKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения дневного счетчика: %v", err)
	}

	// Лимиты (увеличены для тестирования)
	var minuteLimit, dayLimit int
	if limitType == "public" {
		minuteLimit = 2000 // Увеличено для публичных endpoints
		dayLimit = 20000   // Увеличено для публичных endpoints
	} else {
		minuteLimit = 1000 // Увеличено для всех endpoints
		dayLimit = 10000   // Увеличено для всех endpoints
	}

	// Проверяем лимиты
	allowed := true
	message := ""

	if minuteCount >= minuteLimit {
		allowed = false
		message = fmt.Sprintf("Превышен минутный лимит: %d/%d", minuteCount, minuteLimit)
	} else if isGetMethod && dayCount >= dayLimit {
		allowed = false
		message = fmt.Sprintf("Превышен дневной лимит: %d/%d", dayCount, dayLimit)
	}

	// Если запрос разрешен, обновляем счетчики
	if allowed {
		err = s.incrementCounters(apiKey, limitType, isGetMethod)
		if err != nil {
			log.Printf("Ошибка обновления счетчиков: %v", err)
		}
	}

	return &RateLimitCheck{
		Allowed:     allowed,
		MinuteLimit: minuteLimit,
		DayLimit:    dayLimit,
		MinuteUsed:  minuteCount,
		DayUsed:     dayCount,
		Message:     message,
	}, nil
}

// getCounter получает счетчик из Redis
func (s *RedisRateLimitService) getCounter(key string) (int, error) {
	val, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil // Ключ не существует
	}
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга счетчика: %v", err)
	}

	return count, nil
}

// incrementCounters увеличивает счетчики
func (s *RedisRateLimitService) incrementCounters(apiKey, limitType string, isGetMethod bool) error {
	now := time.Now()

	// Минутный счетчик
	minuteKey := fmt.Sprintf("rate_limit:%s:%s:minute", apiKey, limitType)
	minuteExpiry := time.Until(now.Truncate(time.Minute).Add(time.Minute))

	// Дневной счетчик
	dayKey := fmt.Sprintf("rate_limit:%s:%s:day", apiKey, limitType)
	dayExpiry := time.Until(now.Truncate(24 * time.Hour).Add(24 * time.Hour))

	// Транзакция для атомарного обновления
	pipe := s.client.Pipeline()

	// Увеличиваем минутный счетчик
	pipe.Incr(s.ctx, minuteKey)
	pipe.Expire(s.ctx, minuteKey, minuteExpiry)

	// Увеличиваем дневный счетчик только для GET методов
	if isGetMethod {
		pipe.Incr(s.ctx, dayKey)
		pipe.Expire(s.ctx, dayKey, dayExpiry)
	}

	_, err := pipe.Exec(s.ctx)
	return err
}

// GetAPIKeyInfo получает информацию об API ключе
func (s *RedisRateLimitService) GetAPIKeyInfo(apiKey string) (*APIKeyInfo, error) {
	// Получаем статистику по всем эндпоинтам
	endpoints := []string{"all", "public"}

	info := &APIKeyInfo{
		APIKey:        apiKey,
		Endpoints:     make(map[string]*EndpointStats),
		TotalRequests: 0,
		LastRequest:   time.Time{},
	}

	for _, endpoint := range endpoints {
		minuteKey := fmt.Sprintf("rate_limit:%s:%s:minute", apiKey, endpoint)
		dayKey := fmt.Sprintf("rate_limit:%s:%s:day", apiKey, endpoint)

		minuteCount, _ := s.getCounter(minuteKey)
		dayCount, _ := s.getCounter(dayKey)

		// Получаем время последнего запроса
		lastRequestKey := fmt.Sprintf("rate_limit:%s:%s:last_request", apiKey, endpoint)
		lastRequestStr, err := s.client.Get(s.ctx, lastRequestKey).Result()
		var lastRequest time.Time
		if err == nil {
			lastRequest, _ = time.Parse(time.RFC3339, lastRequestStr)
		}

		info.Endpoints[endpoint] = &EndpointStats{
			MinuteCount: minuteCount,
			DayCount:    dayCount,
			LastRequest: lastRequest,
		}

		info.TotalRequests += dayCount
		if lastRequest.After(info.LastRequest) {
			info.LastRequest = lastRequest
		}
	}

	return info, nil
}

// SearchAPIKeys поиск API ключей по паттерну
func (s *RedisRateLimitService) SearchAPIKeys(pattern string) ([]string, error) {
	// Используем SCAN для поиска ключей
	var allKeys []string
	var cursor uint64
	var err error

	searchPattern := fmt.Sprintf("rate_limit:%s:*", pattern)

	for {
		var keys []string
		var scanResult []string
		scanResult, cursor, err = s.client.Scan(s.ctx, cursor, searchPattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = scanResult

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}
	}

	// Извлекаем уникальные API ключи
	apiKeys := make(map[string]bool)
	for _, key := range allKeys {
		// Формат ключа: rate_limit:API_KEY:endpoint:type
		parts := strings.Split(key, ":")
		if len(parts) >= 2 {
			apiKeys[parts[1]] = true
		}
	}

	// Преобразуем map в slice
	result := make([]string, 0, len(apiKeys))
	for apiKey := range apiKeys {
		result = append(result, apiKey)
	}

	return result, nil
}

// GetRateLimitStats получает статистику по rate limiting
func (s *RedisRateLimitService) GetRateLimitStats() (*RateLimitStats, error) {
	// Получаем общую статистику
	info := s.client.Info(s.ctx).Val()

	// Подсчитываем количество активных ключей
	keys, err := s.client.Keys(s.ctx, "rate_limit:*").Result()
	if err != nil {
		return nil, err
	}

	// Группируем по API ключам
	apiKeyCount := make(map[string]int)
	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) >= 2 {
			apiKey := parts[1]
			apiKeyCount[apiKey]++
		}
	}

	stats := &RateLimitStats{
		TotalAPIKeys:  len(apiKeyCount),
		TotalKeys:     len(keys),
		RedisInfo:     info,
		ActiveAPIKeys: make([]string, 0, len(apiKeyCount)),
	}

	for apiKey := range apiKeyCount {
		stats.ActiveAPIKeys = append(stats.ActiveAPIKeys, apiKey)
	}

	return stats, nil
}

// Close закрывает соединение с Redis
func (s *RedisRateLimitService) Close() error {
	return s.client.Close()
}

// ResetCounters сбрасывает счетчики для API ключа
func (s *RedisRateLimitService) ResetCounters(apiKey, resetType string) error {
	endpoints := []string{"all", "public"}

	for _, endpoint := range endpoints {
		if resetType != "all" && resetType != endpoint {
			continue
		}

		// Удаляем ключи счетчиков
		minuteKey := fmt.Sprintf("rate_limit:%s:%s:minute", apiKey, endpoint)
		dayKey := fmt.Sprintf("rate_limit:%s:%s:day", apiKey, endpoint)
		lastRequestKey := fmt.Sprintf("rate_limit:%s:%s:last_request", apiKey, endpoint)

		pipe := s.client.Pipeline()
		pipe.Del(s.ctx, minuteKey)
		pipe.Del(s.ctx, dayKey)
		pipe.Del(s.ctx, lastRequestKey)

		_, err := pipe.Exec(s.ctx)
		if err != nil {
			return fmt.Errorf("ошибка сброса счетчиков для %s: %v", endpoint, err)
		}
	}

	return nil
}
