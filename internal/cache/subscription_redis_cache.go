package cache

import (
	"context"
	"encoding/json"
	"time"

	"effective-project/internal/models"

	"github.com/redis/go-redis/v9"
)

type SubscriptionRedisCache struct {
	rdb *redis.Client
}

func NewSubscriptionRedisCache(rdb *redis.Client) *SubscriptionRedisCache {
	return &SubscriptionRedisCache{
		rdb: rdb,
	}
}

func (c *SubscriptionRedisCache) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	key := "subscription:id:" + id

	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var sub models.Subscription
	if err := json.Unmarshal([]byte(val), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

func (c *SubscriptionRedisCache) Set(
	ctx context.Context,
	sub *models.Subscription,
	ttl time.Duration,
) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	id := sub.ID.String()

	return c.rdb.Set(ctx, "subscription:id:"+id, data, ttl).Err()
}

func (c *SubscriptionRedisCache) DeleteByID(ctx context.Context, id string) error {
	return c.rdb.Del(ctx, "subscription:id:"+id).Err()
}
