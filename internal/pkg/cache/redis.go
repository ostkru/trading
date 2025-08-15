package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // нет пароля
		DB:       0,  // используем базу данных 0
		PoolSize: 10, // размер пула соединений
	})

	ctx := context.Background()

	// Проверяем соединение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Ошибка подключения к Redis: %v\n", err)
		return nil
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

// Set устанавливает значение в кэш
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %v", err)
	}

	return r.client.Set(r.ctx, key, jsonValue, expiration).Err()
}

// Get получает значение из кэша
func (r *RedisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete удаляет ключ из кэша
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists проверяет существование ключа
func (r *RedisCache) Exists(key string) bool {
	exists, _ := r.client.Exists(r.ctx, key).Result()
	return exists > 0
}

// Flush очищает весь кэш
func (r *RedisCache) Flush() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close закрывает соединение
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetStats возвращает статистику кэша
func (r *RedisCache) GetStats() map[string]interface{} {
	info := r.client.Info(r.ctx).Val()

	stats := make(map[string]interface{})
	stats["connected_clients"] = r.client.PoolStats()
	stats["info"] = info

	return stats
}
