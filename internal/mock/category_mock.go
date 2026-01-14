package mock

import (
	"context"
	"time"

	"effective-project/internal/models"

	"github.com/google/uuid"
)

// MockCategoryRepository is a test mock for repository.CategoryRepository
type MockCategoryRepository struct {
	CreateFn  func(c *models.Category) error
	ListFn    func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.Category, error)
	GetByIDFn func(id string) (*models.Category, error)
	UpdateFn  func(c *models.Category) error
	DeleteFn  func(id string) error
}

func (m *MockCategoryRepository) Create(c *models.Category) error {
	if m.CreateFn != nil {
		return m.CreateFn(c)
	}
	return nil
}

func (m *MockCategoryRepository) List(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.Category, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, limit, lastCreatedAt, lastID)
	}
	return nil, nil
}

func (m *MockCategoryRepository) GetByID(id string) (*models.Category, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockCategoryRepository) Update(c *models.Category) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(c)
	}
	return nil
}

func (m *MockCategoryRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

// MockCategoryCache is a mock for cache.CategoryCache
type MockCategoryCache struct {
	GetByIDFn func(ctx context.Context, id string) (*models.Category, error)
	SetFn     func(ctx context.Context, category *models.Category, ttl time.Duration) error
	DeleteFn  func(ctx context.Context, id string) error
}

func (m *MockCategoryCache) GetByID(ctx context.Context, id string) (*models.Category, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockCategoryCache) Set(ctx context.Context, category *models.Category, ttl time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, category, ttl)
	}
	return nil
}

func (m *MockCategoryCache) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
