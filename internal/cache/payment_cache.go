package cache

import (
	"context"
	"time"

	"effective-project/internal/models"
)

type PaymentCache interface {
	GetByID(ctx context.Context, id string) (*models.Payment, error)
	Set(ctx context.Context, payment *models.Payment, ttl time.Duration) error
	Delete(ctx context.Context, id string) error
}
