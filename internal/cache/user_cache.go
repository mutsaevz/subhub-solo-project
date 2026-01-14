package cache

import (
	"context"
	"effective-project/internal/models"
	"time"
)

type UserCache interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	Set(ctx context.Context, user *models.User, ttl time.Duration) error

	DeleteByID(ctx context.Context, id string) error
}
