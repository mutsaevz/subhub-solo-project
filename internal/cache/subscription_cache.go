package cache

import (
	"context"
	"time"

	"effective-project/internal/models"
)

type SubscriptionCache interface {
	GetByID(ctx context.Context, id string) (*models.Subscription, error)
	Set(ctx context.Context, sub *models.Subscription, ttl time.Duration) error
	DeleteByID(ctx context.Context, id string) error
}
