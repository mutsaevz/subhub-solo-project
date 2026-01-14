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

func TestUserService_Create(t *testing.T) {
	called := false
	repo := &mock.MockUserRepository{
		CreateFn: func(u *models.User) error {
			called = true
			return nil
		},
	}

	cache := &mock.MockUserCache{}
	svc := service.NewUserService(repo, cache, nil)

	req := &dto.UserCreateRequest{
		Email:    "test@example.com",
		Password: "pass",
	}

	created, err := svc.Create(req)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, req.Email, created.Email)
}

func TestUserService_GetByID_Success(t *testing.T) {
	user := &models.User{Base: models.Base{ID: uuid.New()}, Email: "a@mail.com"}

	repo := &mock.MockUserRepository{
		GetByIDFn: func(id string) (*models.User, error) {
			return user, nil
		},
	}
	cache := &mock.MockUserCache{}
	svc := service.NewUserService(repo, cache, nil)

	got, err := svc.GetByID(user.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, user, got)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := &mock.MockUserRepository{
		GetByIDFn: func(id string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewUserService(repo, nil, nil)

	got, err := svc.GetByID("1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestUserService_Delete(t *testing.T) {
	deleted := false
	repo := &mock.MockUserRepository{
		GetByIDFn: func(id string) (*models.User, error) {
			return &models.User{Base: models.Base{ID: uuid.New()}}, nil
		},
		DeleteFn: func(id string) error {
			deleted = true
			return nil
		},
	}

	cacheDeleted := false
	cache := &mock.MockUserCache{
		DeleteByIDFn: func(ctx context.Context, id string) error {
			cacheDeleted = true
			return nil
		},
	}

	svc := service.NewUserService(repo, cache, nil)
	err := svc.Delete("1")

	assert.NoError(t, err)
	assert.True(t, deleted)
	assert.True(t, cacheDeleted)
}

func TestUserService_List(t *testing.T) {
	repo := &mock.MockUserRepository{
		ListFn: func(ctx context.Context, limit int, lastCreatedAt *time.Time, lastID *uuid.UUID) ([]models.User, error) {
			return []models.User{{Base: models.Base{ID: uuid.New()}}, {Base: models.Base{ID: uuid.New()}}}, nil
		},
	}

	svc := service.NewUserService(repo, nil, nil)
	list, err := svc.List(context.Background(), 10, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
