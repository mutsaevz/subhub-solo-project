package mock

import (
	"context"
	"time"

	"effective-project/internal/models"

	"github.com/google/uuid"
)

// MockServiceRepository is a test mock for repository.ServiceRepository
type MockServiceRepository struct {
	CreateFn  func(svc *models.Service) error
	ListFn    func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.Service, error)
	GetByIDFn func(id string) (*models.Service, error)
	UpdateFn  func(svc *models.Service) error
	DeleteFn  func(id string) error
}

func (m *MockServiceRepository) Create(svc *models.Service) error {
	if m.CreateFn != nil {
		return m.CreateFn(svc)
	}
	return nil
}

func (m *MockServiceRepository) List(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.Service, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, limit, lastCreatedAt, lastID)
	}
	return nil, nil
}

func (m *MockServiceRepository) GetByID(id string) (*models.Service, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockServiceRepository) Update(svc *models.Service) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(svc)
	}
	return nil
}

func (m *MockServiceRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

// MockCache is a simple mock for cache.Cache
type MockCache struct {
	SetFn    func(ctx context.Context, key string, value any) error
	GetFn    func(ctx context.Context, key string, dest any) (bool, error)
	DeleteFn func(ctx context.Context, key string) error
}

func (m *MockCache) Set(ctx context.Context, key string, value any) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, key, value)
	}
	return nil
}

func (m *MockCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, key, dest)
	}
	return false, nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, key)
	}
	return nil
}
