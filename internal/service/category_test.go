package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"effective-project/internal/dto"
	"effective-project/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mocks ----

type categoryRepoMock struct{ mock.Mock }

func (m *categoryRepoMock) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *categoryRepoMock) List(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.Category, error) {
	args := m.Called(ctx, limit, lastCreatedAt, lastID)
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *categoryRepoMock) GetByID(id string) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *categoryRepoMock) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *categoryRepoMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type categoryCacheMock struct{ mock.Mock }

func (m *categoryCacheMock) GetByID(ctx context.Context, id string) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *categoryCacheMock) Set(ctx context.Context, categories *models.Category, ttl time.Duration) error {
	args := m.Called(ctx, categories, ttl)
	return args.Error(0)
}

func (m *categoryCacheMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newLoggerr() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// ---- Tests ----

func TestCategoryService_Create(t *testing.T) {
	repo := new(categoryRepoMock)
	cache := new(categoryCacheMock)
	logger := newLoggerr()
	service := NewCategoryService(repo, cache, logger)

	req := &dto.CategoryCreateRequest{Name: "Books"}

	repo.On("Create", mock.MatchedBy(func(c *models.Category) bool { return c.Name == req.Name })).Return(nil)
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	category, err := service.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, req.Name, category.Name)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	repo := new(categoryRepoMock)
	cache := new(categoryCacheMock)
	logger := newLoggerr()
	service := NewCategoryService(repo, cache, logger)

	id := uuid.New()
	expected := &models.Category{Base: models.Base{ID: id}, Name: "Electronics"}

	cache.On("GetByID", mock.Anything, id.String()).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", id.String()).Return(expected, nil)
	cache.On("Set", mock.Anything, expected, mock.Anything).Return(nil)

	category, err := service.GetByID(id.String())

	assert.NoError(t, err)
	assert.Equal(t, expected, category)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	repo := new(categoryRepoMock)
	cache := new(categoryCacheMock)
	logger := newLoggerr()
	service := NewCategoryService(repo, cache, logger)

	id := uuid.New()

	cache.On("GetByID", mock.Anything, id.String()).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", id.String()).Return(nil, errors.New("not found"))

	category, err := service.GetByID(id.String())

	assert.Error(t, err)
	assert.Nil(t, category)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	repo := new(categoryRepoMock)
	cache := new(categoryCacheMock)
	logger := newLoggerr()
	service := NewCategoryService(repo, cache, logger)

	id := uuid.New()

	repo.On("Delete", id.String()).Return(nil)
	cache.On("Delete", mock.Anything, id.String()).Return(nil)

	err := service.Delete(id.String())

	assert.NoError(t, err)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCategoryService_List(t *testing.T) {
	repo := new(categoryRepoMock)
	cache := new(categoryCacheMock)
	logger := newLoggerr()
	service := NewCategoryService(repo, cache, logger)

	expected := []models.Category{{Base: models.Base{ID: uuid.New()}, Name: "A"}, {Base: models.Base{ID: uuid.New()}, Name: "B"}}

	repo.On("List", mock.Anything, 10, (*time.Time)(nil), (*uuid.UUID)(nil)).Return(expected, nil)

	list, err := service.List(context.Background(), 10, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, expected, list)

	repo.AssertExpectations(t)
}
