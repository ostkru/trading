package ratelimit

import "time"

// RateLimit представляет запись о лимите запросов
type RateLimit struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	APIKey          string    `json:"api_key"`
	Endpoint        string    `json:"endpoint"`
	RequestCount    int       `json:"request_count"`
	MinuteCount     int       `json:"minute_count"`
	DayCount        int       `json:"day_count"`
	LastRequestTime time.Time `json:"last_request_time"`
	MinuteStart     time.Time `json:"minute_start"`
	DayStart        time.Time `json:"day_start"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RateLimitConfig представляет конфигурацию лимитов
type RateLimitConfig struct {
	MinuteLimit int `json:"minute_limit"`
	DayLimit    int `json:"day_limit"`
}

// RateLimitCheck представляет результат проверки лимита
type RateLimitCheck struct {
	Allowed     bool   `json:"allowed"`
	MinuteLimit int    `json:"minute_limit"`
	DayLimit    int    `json:"day_limit"`
	MinuteUsed  int    `json:"minute_used"`
	DayUsed     int    `json:"day_used"`
	Message     string `json:"message,omitempty"`
}

// APIKeyInfo представляет информацию об API ключе
type APIKeyInfo struct {
	APIKey        string                    `json:"api_key"`
	Endpoints     map[string]*EndpointStats `json:"endpoints"`
	TotalRequests int                       `json:"total_requests"`
	LastRequest   time.Time                 `json:"last_request"`
}

// EndpointStats представляет статистику по эндпоинту
type EndpointStats struct {
	MinuteCount int       `json:"minute_count"`
	DayCount    int       `json:"day_count"`
	LastRequest time.Time `json:"last_request"`
}

// RateLimitStats представляет общую статистику rate limiting
type RateLimitStats struct {
	TotalAPIKeys  int      `json:"total_api_keys"`
	TotalKeys     int      `json:"total_keys"`
	RedisInfo     string   `json:"redis_info"`
	ActiveAPIKeys []string `json:"active_api_keys"`
}
