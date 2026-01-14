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

type orderRepoMock struct{ mock.Mock }

func (m *orderRepoMock) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *orderRepoMock) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *orderRepoMock) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

type orderCacheMock struct{ mock.Mock }

func (m *orderCacheMock) GetByID(ctx context.Context, id string) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *orderCacheMock) Set(ctx context.Context, order *models.Order, ttl time.Duration) error {
	args := m.Called(ctx, order, ttl)
	return args.Error(0)
}

func (m *orderCacheMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---- Tests ----

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestOrderService_Create(t *testing.T) {
	repo := new(orderRepoMock)
	cache := new(orderCacheMock)
	logger := newLogger()
	service := NewOrderService(repo, cache, logger)

	req := dto.OrderCreateRequest{UserID: uuid.New(), ServiceID: uuid.New(), IsPaid: false}

	repo.On("Create", mock.MatchedBy(func(o *models.Order) bool {
		return o.UserID == req.UserID && o.ServiceID == req.ServiceID && o.IsPaid == req.IsPaid
	})).Return(nil)
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	order, err := service.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, req.UserID, order.UserID)
	assert.Equal(t, req.ServiceID, order.ServiceID)
	assert.Equal(t, req.IsPaid, order.IsPaid)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestOrderService_GetByID_Success(t *testing.T) {
	repo := new(orderRepoMock)
	cache := new(orderCacheMock)
	logger := newLogger()
	service := NewOrderService(repo, cache, logger)

	id := uuid.New()
	expected := &models.Order{Base: models.Base{ID: id}, UserID: uuid.New(), ServiceID: uuid.New(), IsPaid: true}

	cache.On("GetByID", mock.Anything, id.String()).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", id.String()).Return(expected, nil)
	cache.On("Set", mock.Anything, expected, mock.Anything).Return(nil)

	order, err := service.GetByID(id.String())

	assert.NoError(t, err)
	assert.Equal(t, expected, order)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestOrderService_GetByID_NotFound(t *testing.T) {
	repo := new(orderRepoMock)
	cache := new(orderCacheMock)
	logger := newLogger()
	service := NewOrderService(repo, cache, logger)

	id := uuid.New()

	cache.On("GetByID", mock.Anything, id.String()).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", id.String()).Return(nil, errors.New("not found"))

	order, err := service.GetByID(id.String())

	assert.Error(t, err)
	assert.Nil(t, order)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestOrderService_Update(t *testing.T) {
	repo := new(orderRepoMock)
	cache := new(orderCacheMock)
	logger := newLogger()
	service := NewOrderService(repo, cache, logger)

	id := uuid.New()
	existing := &models.Order{Base: models.Base{ID: id}, UserID: uuid.New(), ServiceID: uuid.New(), IsPaid: false}

	repo.On("GetByID", id.String()).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(o *models.Order) bool { return o.IsPaid && o.Base.ID == id })).Return(nil)
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	isPaid := true
	updated, err := service.Update(id.String(), dto.OrderUpdateRequest{IsPaid: &isPaid})

	assert.NoError(t, err)
	assert.True(t, updated.IsPaid)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
