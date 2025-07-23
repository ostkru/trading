package ratelimit

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CheckRateLimit проверяет лимиты для пользователя и эндпоинта
func (s *Service) CheckRateLimit(userID int64, apiKey, endpoint string, isGetMethod bool) (*RateLimitCheck, error) {
	// Определяем тип лимита на основе эндпоинта
	limitType := "all" // по умолчанию для всех методов
	
	// Для публичных методов используем отдельный лимит
	if endpoint == "public" {
		limitType = "public"
	}
	
	// Получаем текущую запись о лимитах
	rateLimit, err := s.getOrCreateRateLimit(userID, apiKey, limitType)
	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	
	// Проверяем минутные лимиты (работают для всех методов)
	minuteLimit := 60
	minuteUsed := rateLimit.MinuteCount
	
	// Сброс минутного счетчика если прошла минута
	if now.Sub(rateLimit.MinuteStart) >= time.Minute {
		minuteUsed = 0
		rateLimit.MinuteStart = now
	}
	
	// Проверяем дневные лимиты (только для GET методов)
	dayLimit := 1000
	dayUsed := rateLimit.DayCount
	
	// Сброс дневного счетчика если прошел день
	if now.Sub(rateLimit.DayStart) >= 24*time.Hour {
		dayUsed = 0
		rateLimit.DayStart = now
	}
	
	// Проверяем лимиты
	allowed := true
	message := ""
	
	if minuteUsed >= minuteLimit {
		allowed = false
		message = fmt.Sprintf("Превышен минутный лимит: %d/%d", minuteUsed, minuteLimit)
	} else if isGetMethod && dayUsed >= dayLimit {
		allowed = false
		message = fmt.Sprintf("Превышен дневной лимит: %d/%d", dayUsed, dayLimit)
	}
	
	// Если запрос разрешен, обновляем счетчики
	if allowed {
		// Для минутных лимитов увеличиваем всегда
		newMinuteCount := minuteUsed + 1
		
		// Для дневных лимитов увеличиваем только для GET методов
		newDayCount := dayUsed
		if isGetMethod {
			newDayCount = dayUsed + 1
		}
		
		err = s.updateRateLimit(userID, apiKey, limitType, newMinuteCount, newDayCount, now)
		if err != nil {
			log.Printf("Ошибка обновления лимитов: %v", err)
		}
	}
	
	return &RateLimitCheck{
		Allowed:     allowed,
		MinuteLimit: minuteLimit,
		DayLimit:    dayLimit,
		MinuteUsed:  minuteUsed,
		DayUsed:     dayUsed,
		Message:     message,
	}, nil
}

// getOrCreateRateLimit получает или создает запись о лимитах
func (s *Service) getOrCreateRateLimit(userID int64, apiKey, limitType string) (*RateLimit, error) {
	query := `SELECT id, user_id, api_key, endpoint, request_count, minute_count, day_count, 
	          last_request_time, minute_start, day_start, created_at, updated_at 
	          FROM api_rate_limits 
	          WHERE user_id = ? AND api_key = ? AND endpoint = ?`
	
	var rateLimit RateLimit
	err := s.db.QueryRow(query, userID, apiKey, limitType).Scan(
		&rateLimit.ID,
		&rateLimit.UserID,
		&rateLimit.APIKey,
		&rateLimit.Endpoint,
		&rateLimit.RequestCount,
		&rateLimit.MinuteCount,
		&rateLimit.DayCount,
		&rateLimit.LastRequestTime,
		&rateLimit.MinuteStart,
		&rateLimit.DayStart,
		&rateLimit.CreatedAt,
		&rateLimit.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		// Создаем новую запись
		now := time.Now()
		insertQuery := `INSERT INTO api_rate_limits (user_id, api_key, endpoint, request_count, minute_count, day_count, 
		               last_request_time, minute_start, day_start, created_at, updated_at) 
		               VALUES (?, ?, ?, 0, 0, 0, ?, ?, ?, ?, ?)`
		
		result, err := s.db.Exec(insertQuery, userID, apiKey, limitType, now, now, now, now, now)
		if err != nil {
			return nil, err
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		
		return &RateLimit{
			ID:              id,
			UserID:          userID,
			APIKey:          apiKey,
			Endpoint:        limitType,
			RequestCount:    0,
			MinuteCount:     0,
			DayCount:        0,
			LastRequestTime: now,
			MinuteStart:     now,
			DayStart:        now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &rateLimit, nil
}

// updateRateLimit обновляет счетчики лимитов
func (s *Service) updateRateLimit(userID int64, apiKey, limitType string, minuteCount, dayCount int, now time.Time) error {
	query := `UPDATE api_rate_limits 
	          SET minute_count = ?, day_count = ?, last_request_time = ?, updated_at = ? 
	          WHERE user_id = ? AND api_key = ? AND endpoint = ?`
	
	_, err := s.db.Exec(query, minuteCount, dayCount, now, now, userID, apiKey, limitType)
	return err
}

// GetRateLimitStats получает статистику лимитов для пользователя
func (s *Service) GetRateLimitStats(userID int64, apiKey string) ([]*RateLimit, error) {
	query := `SELECT id, user_id, api_key, endpoint, request_count, minute_count, day_count, 
	          last_request_time, minute_start, day_start, created_at, updated_at 
	          FROM api_rate_limits 
	          WHERE user_id = ? AND api_key = ?`
	
	rows, err := s.db.Query(query, userID, apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rateLimits []*RateLimit
	for rows.Next() {
		var rateLimit RateLimit
		err := rows.Scan(
			&rateLimit.ID,
			&rateLimit.UserID,
			&rateLimit.APIKey,
			&rateLimit.Endpoint,
			&rateLimit.RequestCount,
			&rateLimit.MinuteCount,
			&rateLimit.DayCount,
			&rateLimit.LastRequestTime,
			&rateLimit.MinuteStart,
			&rateLimit.DayStart,
			&rateLimit.CreatedAt,
			&rateLimit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rateLimits = append(rateLimits, &rateLimit)
	}
	
	return rateLimits, nil
}

// ResetRateLimits сбрасывает лимиты для пользователя (для админских целей)
func (s *Service) ResetRateLimits(userID int64, apiKey string) error {
	query := `UPDATE api_rate_limits 
	          SET minute_count = 0, day_count = 0, last_request_time = NOW(), updated_at = NOW() 
	          WHERE user_id = ? AND api_key = ?`
	
	_, err := s.db.Exec(query, userID, apiKey)
	return err
} 