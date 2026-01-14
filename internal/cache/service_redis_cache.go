package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// ServiceCache — базовый сервис для работы с Redis-кешем
type ServiceCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewServiceCache — конструктор кеш-сервиса
func NewServiceCache(client *redis.Client, ttl time.Duration) *ServiceCache {
	return &ServiceCache{
		client: client,
		ttl:    ttl,
	}
}

// Set — сохранить данные в Redis по ключу
func (c *ServiceCache) Set(ctx context.Context, key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bytes, c.ttl).Err()
}

// Get — получить данные из Redis по ключу
func (c *ServiceCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}

	return true, nil
}

// Delete — удалить ключ из Redis
func (c *ServiceCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
