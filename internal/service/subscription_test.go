package service_test

import (
	"context"
	"effective-project/internal/dto"
	"effective-project/internal/mock"
	"effective-project/internal/models"
	service "effective-project/internal/service"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---- Mocks ----

// use internal/mock.MockSubscriptionRepository

// ---- Tests ----
func timePtr(t time.Time) *time.Time {
	return &t
}
func TestSubscriptionService_Create(t *testing.T) {
	called := false
	repo := &mock.MockSubscriptionRepository{
		CreateFn: func(s *models.Subscription) error {
			called = true
			return nil
		},
	}

	svc := service.NewSubscriptionService(repo, nil, nil, nil, nil)

	req := &dto.SubscriptionCreateRequest{
		UserID:    uuid.New(),
		ServiceID: uuid.New(),
		StartDate: time.Now(),
		EndDate:   timePtr(time.Now().Add(30 * 24 * time.Hour)),
		Price:     100,
	}

	created, err := svc.Create(req)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, req.UserID, created.UserID)
}

func TestSubscriptionService_GetByID_Success(t *testing.T) {
	repo := &mock.MockSubscriptionRepository{
		GetByIDFn: func(id string) (*dto.SubscriptionResponse, error) {
			return &dto.SubscriptionResponse{
				ID:        uuid.New(),
				ServiceID: uuid.New(),
				StartDate: time.Now(),
				EndDate:   nil,
				Price:     100,
				CreatedAt: time.Now(),
			}, nil
		},
		GetModelByIDFn: func(id string) (*models.Subscription, error) {
			return &models.Subscription{Base: models.Base{ID: uuid.New()}, ServiceID: uuid.New()}, nil
		},
	}

	svc := service.NewSubscriptionService(repo, nil, nil, nil, nil)

	id := uuid.New()
	sub, err := svc.GetByID(id.String())

	assert.NoError(t, err)
	assert.NotNil(t, sub)
}

func TestSubscriptionService_GetByID_NotFound(t *testing.T) {
	repo := &mock.MockSubscriptionRepository{
		GetByIDFn: func(id string) (*dto.SubscriptionResponse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewSubscriptionService(repo, nil, nil, nil, nil)

	id := uuid.New()
	sub, err := svc.GetByID(id.String())

	assert.Error(t, err)
	assert.Nil(t, sub)
}

func TestSubscriptionService_Delete(t *testing.T) {
	called := false
	repo := &mock.MockSubscriptionRepository{
		DeleteFn: func(id string) error {
			called = true
			return nil
		},
	}

	svc := service.NewSubscriptionService(repo, nil, nil, nil, nil)

	id := uuid.New()
	err := svc.Delete(id.String())

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSubscriptionService_List(t *testing.T) {
	repo := &mock.MockSubscriptionRepository{
		ListFn: func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]dto.SubscriptionResponse, error) {
			return []dto.SubscriptionResponse{{ID: uuid.New()}, {ID: uuid.New()}}, nil
		},
	}

	svc := service.NewSubscriptionService(repo, nil, nil, nil, nil)

	list, err := svc.List(context.Background(), 10, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
