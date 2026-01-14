package cache

import (
	"context"
	"effective-project/internal/models"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type OrdersRedisCache struct {
	client *redis.Client
	prefix string
}

func NewOrdersRedisCache(client *redis.Client) *OrdersRedisCache {
	return &OrdersRedisCache{
		client: client,
		prefix: "order:",
	}
}

func (c *OrdersRedisCache) key(id string) string {
	return c.prefix + id
}

func (c *OrdersRedisCache) GetByID(ctx context.Context, id string) (*models.Order, error) {
	data, err := c.client.Get(ctx, c.key(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *OrdersRedisCache) Set(ctx context.Context, order *models.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	id := order.ID.String()

	return c.client.Set(ctx, c.key(id), data, ttl).Err()
}

func (c *OrdersRedisCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
