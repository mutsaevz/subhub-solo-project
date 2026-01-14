package cache

import (
	"context"
	"effective-project/internal/models"
	"time"
)

type CategoryCache interface {
	GetByID(ctx context.Context, id string) (*models.Category, error)

	Set(ctx context.Context, category *models.Category, ttl time.Duration) error

	Delete(ctx context.Context, id string) error
}
