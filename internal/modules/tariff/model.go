package tariff

import "time"

// Tariff представляет тарифный план
type Tariff struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MinuteLimit int       `json:"minute_limit"`
	DayLimit    int       `json:"day_limit"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTariffRequest представляет запрос на создание тарифа
type CreateTariffRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	MinuteLimit int    `json:"minute_limit" binding:"required,min=1"`
	DayLimit    int    `json:"day_limit" binding:"required,min=1"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// UpdateTariffRequest представляет запрос на обновление тарифа
type UpdateTariffRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	MinuteLimit *int    `json:"minute_limit,omitempty" binding:"omitempty,min=1"`
	DayLimit    *int    `json:"day_limit,omitempty" binding:"omitempty,min=1"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// TariffListResponse представляет ответ со списком тарифов
type TariffListResponse struct {
	Tariffs []Tariff `json:"tariffs"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

// UserTariffInfo представляет информацию о тарифе пользователя
type UserTariffInfo struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	TariffID    int64  `json:"tariff_id"`
	TariffName  string `json:"tariff_name"`
	MinuteLimit int    `json:"minute_limit"`
	DayLimit    int    `json:"day_limit"`
	IsActive    bool   `json:"is_active"`
}

// TariffUsageStats представляет статистику использования тарифа
type TariffUsageStats struct {
	TariffID       int64   `json:"tariff_id"`
	TariffName     string  `json:"tariff_name"`
	UserCount      int     `json:"user_count"`
	MinuteLimit    int     `json:"minute_limit"`
	DayLimit       int     `json:"day_limit"`
	AvgMinuteUsage float64 `json:"avg_minute_usage"`
	AvgDayUsage    float64 `json:"avg_day_usage"`
	LastActivity   *string `json:"last_activity,omitempty"`
}

// ChangeUserTariffRequest представляет запрос на изменение тарифа пользователя
type ChangeUserTariffRequest struct {
	UserID   int64 `json:"user_id" binding:"required,min=1"`
	TariffID int64 `json:"tariff_id" binding:"required,min=1"`
}

// TariffLimits представляет лимиты тарифа
type TariffLimits struct {
	MinuteLimit int `json:"minute_limit"`
	DayLimit    int `json:"day_limit"`
}
