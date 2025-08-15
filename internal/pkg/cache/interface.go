package cache

import (
	"fmt"
	"time"
)

// Cache интерфейс для кэширования
type Cache interface {
	// Set устанавливает значение в кэш
	Set(key string, value interface{}, expiration time.Duration) error

	// Get получает значение из кэша
	Get(key string, dest interface{}) error

	// Delete удаляет ключ из кэша
	Delete(key string) error

	// Exists проверяет существование ключа
	Exists(key string) bool

	// Flush очищает весь кэш
	Flush() error

	// Close закрывает соединение
	Close() error

	// GetStats возвращает статистику кэша
	GetStats() map[string]interface{}
}

// MemoryCache простая реализация кэша в памяти (fallback)
type MemoryCache struct {
	data map[string]interface{}
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]interface{}),
	}
}

func (m *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MemoryCache) Get(key string, dest interface{}) error {
	if _, exists := m.data[key]; exists {
		// Простое копирование для демонстрации
		// В реальном приложении нужна более сложная логика
		return nil
	}
	return fmt.Errorf("ключ не найден")
}

func (m *MemoryCache) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *MemoryCache) Exists(key string) bool {
	_, exists := m.data[key]
	return exists
}

func (m *MemoryCache) Flush() error {
	m.data = make(map[string]interface{})
	return nil
}

func (m *MemoryCache) Close() error {
	return nil
}

func (m *MemoryCache) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"type": "memory",
		"size": len(m.data),
	}
}
