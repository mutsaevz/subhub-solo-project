package mock

import (
	"context"
	"time"

	"effective-project/internal/models"
)

// MockOrderRepository is a test mock for repository.OrderRepository
type MockOrderRepository struct {
	CreateFn  func(order *models.Order) error
	GetByIDFn func(id string) (*models.Order, error)
	UpdateFn  func(order *models.Order) error
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	if m.CreateFn != nil {
		return m.CreateFn(order)
	}
	return nil
}

func (m *MockOrderRepository) GetByID(id string) (*models.Order, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockOrderRepository) Update(order *models.Order) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(order)
	}
	return nil
}

// MockOrderCache is a mock for cache.OrderCache
type MockOrderCache struct {
	GetByIDFn func(ctx context.Context, id string) (*models.Order, error)
	SetFn     func(ctx context.Context, order *models.Order, ttl time.Duration) error
	DeleteFn  func(ctx context.Context, id string) error
}

func (m *MockOrderCache) GetByID(ctx context.Context, id string) (*models.Order, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderCache) Set(ctx context.Context, order *models.Order, ttl time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, order, ttl)
	}
	return nil
}

func (m *MockOrderCache) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
