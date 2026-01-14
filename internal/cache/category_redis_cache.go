package cache

import (
	"context"
	"effective-project/internal/models"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CategoryRedisCache struct {
	client *redis.Client
	prefix string
}

func NewCategoryRedisCache(client *redis.Client) *CategoryRedisCache {
	return &CategoryRedisCache{
		client: client,
		prefix: "category:",
	}
}

func (c *CategoryRedisCache) key(id string) string {
	return c.prefix + id
}

func (c *CategoryRedisCache) GetByID(ctx context.Context, id string) (*models.Category, error) {
	data, err := c.client.Get(ctx, c.key(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var category models.Category
	if err := json.Unmarshal([]byte(data), &category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *CategoryRedisCache) Set(ctx context.Context, category *models.Category, ttl time.Duration) error {
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}

	id := category.ID.String()

	return c.client.Set(ctx, c.key(id), data, ttl).Err()
}

func (c *CategoryRedisCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
