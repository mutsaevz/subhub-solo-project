package cache

import (
	"context"
	"effective-project/internal/models"
	"time"
)

type OrderCache interface {
	GetByID(ctx context.Context, id string) (*models.Order, error)

	Set(ctx context.Context, order *models.Order, ttl time.Duration) error

	Delete(ctx context.Context, id string) error
}
