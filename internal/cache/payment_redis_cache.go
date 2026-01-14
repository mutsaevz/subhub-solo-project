package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"effective-project/internal/models"

	"github.com/redis/go-redis/v9"
)

var ErrPaymentNotFound = errors.New("payment not found")

type PaymentRedisCache struct {
	client *redis.Client
	prefix string
}

func NewPaymentRedisCache(client *redis.Client) *PaymentRedisCache {
	return &PaymentRedisCache{
		client: client,
		prefix: "payment:",
	}
}

func (c *PaymentRedisCache) key(id string) string {
	return fmt.Sprintf("%s%s", c.prefix, id)
}

func (c *PaymentRedisCache) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	val, err := c.client.Get(ctx, c.key(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	var payment models.Payment
	if err := json.Unmarshal([]byte(val), &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (c *PaymentRedisCache) Set(ctx context.Context, payment *models.Payment, ttl time.Duration) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	id := payment.ID.String()

	return c.client.Set(ctx, c.key(id), data, ttl).Err()
}

func (c *PaymentRedisCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
