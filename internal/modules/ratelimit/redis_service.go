package ratelimit

import (
	"context"
	"database/sql"
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
	db     *sql.DB
}

func NewRedisRateLimitService(redisAddr, redisPassword string, redisDB int, db *sql.DB) *RedisRateLimitService {
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
		db:     db,
	}
}

// CheckRateLimit проверяет лимиты для API ключа
func (s *RedisRateLimitService) CheckRateLimit(apiKey, endpoint string, isGetMethod bool) (*RateLimitCheck, error) {
	log.Printf("DEBUG: CheckRateLimit called with apiKey=%s, endpoint=%s", apiKey, endpoint)
	
	// Получаем userID по API ключу
	userID, err := s.getUserIDByAPIKey(apiKey)
	if err != nil {
		// Если пользователь не найден, используем дефолтные лимиты
		log.Printf("Пользователь не найден для API ключа %s, используем дефолтные лимиты", apiKey)
		return s.checkRateLimitWithLimits(apiKey, endpoint, isGetMethod, &RateLimitConfig{
			MinuteLimit: 1000,
			DayLimit:    10000,
		})
	}

	// Получаем лимиты пользователя из тарифа
	minuteLimit, dayLimit, err := s.getUserLimitsFromAPIKey(apiKey)
	if err != nil {
		log.Printf("Ошибка получения лимитов пользователя %d: %v, используем дефолтные лимиты", userID, err)
		minuteLimit = 60
		dayLimit = 1000
	}

	limits := &RateLimitConfig{
		MinuteLimit: minuteLimit,
		DayLimit:    dayLimit,
	}

	return s.checkRateLimitWithLimits(apiKey, endpoint, isGetMethod, limits)
}

// checkRateLimitWithLimits проверяет лимиты с заданными значениями
func (s *RedisRateLimitService) checkRateLimitWithLimits(apiKey, endpoint string, isGetMethod bool, limits *RateLimitConfig) (*RateLimitCheck, error) {
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

	// Используем лимиты из тарифа
	minuteLimit := limits.MinuteLimit
	dayLimit := limits.DayLimit

	// Для публичных endpoints оставляем лимиты как есть (уже настроены в тарифе)
	// Убираем удвоение лимитов, так как тариф уже содержит правильные значения

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

	// Обновляем счетчики для всех запросов (для статистики)
	err = s.incrementCounters(apiKey, limitType, isGetMethod)
	if err != nil {
		log.Printf("Ошибка обновления счетчиков: %v", err)
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

// getUserIDByAPIKey получает ID пользователя по API ключу
func (s *RedisRateLimitService) getUserIDByAPIKey(apiKey string) (int64, error) {
	if s.db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_token = ? AND is_active = TRUE", apiKey).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// getUserTariffLimits получает лимиты тарифа пользователя
func (s *RedisRateLimitService) getUserTariffLimits(userID int64) (*RateLimitConfig, error) {
	if s.db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	var limits RateLimitConfig
	err := s.db.QueryRow(`
		SELECT t.minute_limit, t.day_limit
		FROM users u
		JOIN tariffs t ON u.tariff_id = t.id
		WHERE u.id = ? AND u.is_active = TRUE AND t.is_active = TRUE
	`, userID).Scan(&limits.MinuteLimit, &limits.DayLimit)

	if err != nil {
		return nil, err
	}

	return &limits, nil
}

// getUserLimitsFromAPIKey получает лимиты пользователя по API ключу
func (s *RedisRateLimitService) getUserLimitsFromAPIKey(apiKey string) (int, int, error) {
	if s.db == nil {
		return 0, 0, fmt.Errorf("база данных не инициализирована")
	}

	// Сначала получаем userID по API ключу
	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_token = ? AND is_active = TRUE", apiKey).Scan(&userID)
	if err != nil {
		return 0, 0, err
	}

	// Затем получаем лимиты пользователя
	var dailyLimit int

	// Получаем лимиты из тарифа пользователя
	var minuteLimit int
	query := `
		SELECT t.minute_limit, t.day_limit
		FROM users u
		JOIN tariffs t ON u.tariff_id = t.id
		WHERE u.id = ? AND u.is_active = 1 AND t.is_active = 1
	`
	err = s.db.QueryRow(query, userID).Scan(&minuteLimit, &dailyLimit)
	if err != nil {
		// Если пользователь не найден или тариф не найден, используем дефолтные лимиты
		return 60, 1000, nil
	}

	// Если лимиты не найдены или некорректные, используем дефолтные
	if dailyLimit <= 0 {
		dailyLimit = 1000
	}
	if minuteLimit <= 0 {
		minuteLimit = 60
	}

	// Отладочный лог
	log.Printf("DEBUG: getUserLimitsFromAPIKey apiKey=%s, minuteLimit=%d, dailyLimit=%d", apiKey, minuteLimit, dailyLimit)

	return minuteLimit, dailyLimit, nil
}
