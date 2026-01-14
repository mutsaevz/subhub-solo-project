package mock

import (
	"context"
	"time"

	"effective-project/internal/models"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	CreateFn     func(user *models.User) error
	ListFn       func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.User, error)
	GetByIDFn    func(id string) (*models.User, error)
	GetByEmailFn func(email string) (*models.User, error)
	UpdateFn     func(user *models.User) error
	DeleteFn     func(id string) error
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(user)
	}
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.User, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, limit, lastCreatedAt, lastID)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(email)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(user)
	}
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

// MockUserCache provides a simple mock for cache.UserCache
type MockUserCache struct {
	GetByIDFn    func(ctx context.Context, id string) (*models.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*models.User, error)
	SetFn        func(ctx context.Context, user *models.User, ttl time.Duration) error
	DeleteByIDFn func(ctx context.Context, id string) error
}

func (m *MockUserCache) GetByID(ctx context.Context, id string) (*models.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserCache) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockUserCache) Set(ctx context.Context, user *models.User, ttl time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, user, ttl)
	}
	return nil
}

func (m *MockUserCache) DeleteByID(ctx context.Context, id string) error {
	if m.DeleteByIDFn != nil {
		return m.DeleteByIDFn(ctx, id)
	}
	return nil
}
