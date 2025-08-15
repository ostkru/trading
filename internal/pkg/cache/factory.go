package cache

import (
	"fmt"
	"log"
	"time"
)

// CacheFactory фабрика для создания кэша
type CacheFactory struct{}

// NewCacheFactory создает новую фабрику кэша
func NewCacheFactory() *CacheFactory {
	return &CacheFactory{}
}

// CreateCache создает кэш (Redis или Memory)
func (cf *CacheFactory) CreateCache() Cache {
	// Пытаемся создать Redis кэш
	redisCache := NewRedisCache()
	if redisCache != nil {
		log.Println("✅ Redis кэш успешно инициализирован")
		return redisCache
	}

	// Fallback на Memory кэш
	log.Println("⚠️ Redis недоступен, используется Memory кэш")
	return NewMemoryCache()
}

// CreateProductCache создает кэш для продуктов
func (cf *CacheFactory) CreateProductCache() Cache {
	cache := cf.CreateCache()

	// Устанавливаем TTL для продуктов (1 час)
	// В реальном приложении можно настраивать разные TTL для разных типов данных
	return &ProductCache{
		cache: cache,
		ttl:   1 * time.Hour,
	}
}

// ProductCache специализированный кэш для продуктов
type ProductCache struct {
	cache Cache
	ttl   time.Duration
}

// Set устанавливает значение в кэш
func (pc *ProductCache) Set(key string, value interface{}, expiration time.Duration) error {
	return pc.cache.Set(key, value, expiration)
}

// Get получает значение из кэша
func (pc *ProductCache) Get(key string, dest interface{}) error {
	return pc.cache.Get(key, dest)
}

// Delete удаляет ключ из кэша
func (pc *ProductCache) Delete(key string) error {
	return pc.cache.Delete(key)
}

// Exists проверяет существование ключа
func (pc *ProductCache) Exists(key string) bool {
	return pc.cache.Exists(key)
}

// GetStats возвращает статистику кэша
func (pc *ProductCache) GetStats() map[string]interface{} {
	return pc.cache.GetStats()
}

// SetProduct кэширует продукт
func (pc *ProductCache) SetProduct(productID int64, product interface{}) error {
	key := fmt.Sprintf("product:%d", productID)
	return pc.cache.Set(key, product, pc.ttl)
}

// GetProduct получает продукт из кэша
func (pc *ProductCache) GetProduct(productID int64, dest interface{}) error {
	key := fmt.Sprintf("product:%d", productID)
	return pc.cache.Get(key, dest)
}

// DeleteProduct удаляет продукт из кэша
func (pc *ProductCache) DeleteProduct(productID int64) error {
	key := fmt.Sprintf("product:%d", productID)
	return pc.cache.Delete(key)
}

// SetProductList кэширует список продуктов
func (pc *ProductCache) SetProductList(key string, products interface{}) error {
	cacheKey := fmt.Sprintf("product_list:%s", key)
	return pc.cache.Set(cacheKey, products, pc.ttl)
}

// GetProductList получает список продуктов из кэша
func (pc *ProductCache) GetProductList(key string, dest interface{}) error {
	cacheKey := fmt.Sprintf("product_list:%s", key)
	return pc.cache.Get(cacheKey, dest)
}

// InvalidateProductList инвалидирует кэш списка продуктов
func (pc *ProductCache) InvalidateProductList(key string) error {
	cacheKey := fmt.Sprintf("product_list:%s", key)
	return pc.cache.Delete(cacheKey)
}

// Flush очищает весь кэш продуктов
func (pc *ProductCache) Flush() error {
	return pc.cache.Flush()
}

// Close закрывает кэш
func (pc *ProductCache) Close() error {
	return pc.cache.Close()
}
