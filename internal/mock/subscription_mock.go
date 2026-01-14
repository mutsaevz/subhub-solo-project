package mock

import (
	"context"
	"time"

	"effective-project/internal/dto"
	"effective-project/internal/models"

	"github.com/google/uuid"
)

// MockSubscriptionRepository is a test mock for repository.SubscriptionRepository
type MockSubscriptionRepository struct {
	CreateFn       func(s *models.Subscription) error
	ListFn         func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]dto.SubscriptionResponse, error)
	GetByIDFn      func(id string) (*dto.SubscriptionResponse, error)
	UpdateFn       func(s *models.Subscription) error
	DeleteFn       func(id string) error
	FindForTotalFn func(ctx context.Context, f dto.TotalFilter) ([]dto.SubscriptionRow, error)
	GetModelByIDFn func(id string) (*models.Subscription, error)
}

func (m *MockSubscriptionRepository) Create(s *models.Subscription) error {
	if m.CreateFn != nil {
		return m.CreateFn(s)
	}
	return nil
}

func (m *MockSubscriptionRepository) List(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]dto.SubscriptionResponse, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, limit, lastCreatedAt, lastID)
	}
	return nil, nil
}

func (m *MockSubscriptionRepository) GetByID(id string) (*dto.SubscriptionResponse, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockSubscriptionRepository) Update(s *models.Subscription) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(s)
	}
	return nil
}

func (m *MockSubscriptionRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func (m *MockSubscriptionRepository) FindForTotal(ctx context.Context, f dto.TotalFilter) ([]dto.SubscriptionRow, error) {
	if m.FindForTotalFn != nil {
		return m.FindForTotalFn(ctx, f)
	}
	return nil, nil
}

func (m *MockSubscriptionRepository) GetModelByID(id string) (*models.Subscription, error) {
	if m.GetModelByIDFn != nil {
		return m.GetModelByIDFn(id)
	}
	return nil, nil
}

// MockSubscriptionCache is a mock for cache.SubscriptionCache
type MockSubscriptionCache struct {
	GetByIDFn    func(ctx context.Context, id string) (*models.Subscription, error)
	SetFn        func(ctx context.Context, sub *models.Subscription, ttl time.Duration) error
	DeleteByIDFn func(ctx context.Context, id string) error
}

func (m *MockSubscriptionCache) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSubscriptionCache) Set(ctx context.Context, sub *models.Subscription, ttl time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, sub, ttl)
	}
	return nil
}

func (m *MockSubscriptionCache) DeleteByID(ctx context.Context, id string) error {
	if m.DeleteByIDFn != nil {
		return m.DeleteByIDFn(ctx, id)
	}
	return nil
}
