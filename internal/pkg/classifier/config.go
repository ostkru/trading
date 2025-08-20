package classifier

import (
	"fmt"
	"os"
	"strconv"
)

// Config конфигурация классификатора
type Config struct {
	// API базовый URL
	APIBaseURL string
	// Таймаут HTTP запросов в секундах
	HTTPTimeout int
	// Размер очереди классификации
	QueueSize int
	// Задержка между запросами в миллисекундах
	RequestDelay int
	// Минимальная вероятность для классификации (0.99 = 99%)
	MinConfidence float64
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:    "https://api.ostk.ru",
		HTTPTimeout:   30,
		QueueSize:     1000,
		RequestDelay:  100,
		MinConfidence: 0.99,
	}
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *Config) LoadFromEnv() {
	if apiURL := os.Getenv("OSTK_API_BASE_URL"); apiURL != "" {
		c.APIBaseURL = apiURL
	}

	if timeout := os.Getenv("OSTK_HTTP_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			c.HTTPTimeout = val
		}
	}

	if queueSize := os.Getenv("OSTK_QUEUE_SIZE"); queueSize != "" {
		if val, err := strconv.Atoi(queueSize); err == nil {
			c.QueueSize = val
		}
	}

	if delay := os.Getenv("OSTK_REQUEST_DELAY"); delay != "" {
		if val, err := strconv.Atoi(delay); err == nil {
			c.RequestDelay = val
		}
	}

	if confidence := os.Getenv("OSTK_MIN_CONFIDENCE"); confidence != "" {
		if val, err := strconv.ParseFloat(confidence, 64); err == nil {
			c.MinConfidence = val
		}
	}
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.APIBaseURL == "" {
		return fmt.Errorf("APIBaseURL не может быть пустым")
	}

	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTPTimeout должен быть больше 0")
	}

	if c.QueueSize <= 0 {
		return fmt.Errorf("QueueSize должен быть больше 0")
	}

	if c.RequestDelay < 0 {
		return fmt.Errorf("RequestDelay не может быть отрицательным")
	}

	if c.MinConfidence < 0 || c.MinConfidence > 1 {
		return fmt.Errorf("MinConfidence должен быть в диапазоне [0, 1]")
	}

	return nil
}
