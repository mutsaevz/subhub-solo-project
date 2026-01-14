package cache

import (
	"context"
	"encoding/json"
	"time"

	"effective-project/internal/models"

	"github.com/redis/go-redis/v9"
)

type UserRedisCache struct {
	rdb *redis.Client
}

func NewUserRedisCache(rdb *redis.Client) *UserRedisCache {
	return &UserRedisCache{
		rdb: rdb,
	}
}

func (c *UserRedisCache) GetByID(ctx context.Context, id string) (*models.User, error) {
	key := "user:id:" + id

	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, ErrCacheMiss
	}

	return &user, nil
}

func (c *UserRedisCache) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	key := "user:email:" + email

	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, ErrCacheMiss
	}

	return &user, nil
}

func (c *UserRedisCache) Set(ctx context.Context, user *models.User, ttl time.Duration) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	id := user.ID.String()

	if err := c.rdb.Set(ctx, "user:id:"+id, data, ttl).Err(); err != nil {
		return err
	}
	if err := c.rdb.Set(ctx, "user:email:"+user.Email, data, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (c *UserRedisCache) DeleteByID(ctx context.Context, id string) error {
	return c.rdb.Del(ctx, "user:id:"+id).Err()
}

func (c *UserRedisCache) DeleteByEmail(ctx context.Context, email string) error {
	return c.rdb.Del(ctx, "user:email:"+email).Err()
}
