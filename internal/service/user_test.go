package service

import (
	"context"
	"effective-project/internal/cache"
	"effective-project/internal/dto"
	"effective-project/internal/mock"
	"effective-project/internal/models"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name      string
		req       *dto.UserCreateRequest
		repoErr   error
		wantErr   bool
		checkFunc func(t *testing.T, saved *models.User)
	}{
		{
			name: "success with default role and hashed password",
			req: &dto.UserCreateRequest{
				Email:    "test@mail.com",
				Password: "plain",
			},
			checkFunc: func(t *testing.T, saved *models.User) {
				if saved.Roles != models.RoleUser {
					t.Fatalf("expected role %s, got %s", models.RoleUser, saved.Roles)
				}
				if err := bcrypt.CompareHashAndPassword(
					[]byte(saved.Password),
					[]byte("plain"),
				); err != nil {
					t.Fatalf("password is not hashed")
				}
			},
		},
		{
			name:    "repo error",
			req:     &dto.UserCreateRequest{Email: "err@mail.com", Password: "123"},
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var savedUser *models.User

			repo := &mock.MockUserRepository{
				CreateFn: func(u *models.User) error {
					savedUser = u
					return tt.repoErr
				},
			}
			cache := &mock.MockUserCache{}

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
			svc := NewUserService(repo, cache, logger)
			_, err := svc.Create(tt.req)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, savedUser)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	user := &models.User{Base: models.Base{ID: uuid.New()}, Email: "a@mail.com"}

	tests := []struct {
		name          string
		cacheResult   *models.User
		cacheErr      error
		repoResult    *models.User
		repoErr       error
		expectRepoHit bool
	}{
		{
			name:        "cache hit",
			cacheResult: user,
		},
		{
			name:          "cache miss -> repo",
			cacheErr:      cache.ErrCacheMiss,
			repoResult:    user,
			expectRepoHit: true,
		},
		{
			name:          "repo error",
			cacheErr:      cache.ErrCacheMiss,
			repoErr:       errors.New("db error"),
			expectRepoHit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoCalled := false
			cacheSetCalled := false

			repo := &mock.MockUserRepository{
				GetByIDFn: func(id string) (*models.User, error) {
					repoCalled = true
					return tt.repoResult, tt.repoErr
				},
			}

			cacheMock := &mock.MockUserCache{
				GetByIDFn: func(ctx context.Context, id string) (*models.User, error) {
					return tt.cacheResult, tt.cacheErr
				},
				SetFn: func(ctx context.Context, u *models.User, ttl time.Duration) error {
					cacheSetCalled = true
					return nil
				},
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
			svc := NewUserService(repo, cacheMock, logger)
			_, _ = svc.GetByID("1")

			if tt.expectRepoHit && !repoCalled {
				t.Fatalf("expected repo.GetByID to be called")
			}
			if tt.cacheErr == cache.ErrCacheMiss && tt.repoErr == nil && !cacheSetCalled {
				t.Fatalf("expected cache.Set to be called")
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	user := &models.User{Base: models.Base{ID: uuid.New()}, Email: "old@mail.com"}

	repo := &mock.MockUserRepository{
		GetByIDFn: func(id string) (*models.User, error) {
			return user, nil
		},
		UpdateFn: func(u *models.User) error {
			return nil
		},
	}

	cacheDeleted := false
	cacheMock := &mock.MockUserCache{
		DeleteByIDFn: func(ctx context.Context, id string) error {
			cacheDeleted = true
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	svc := NewUserService(repo, cacheMock, logger)
	newEmail := "new@mail.com"
	_, err := svc.Update(user.ID.String(), &dto.UserUpdateRequest{Email: &newEmail})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cacheDeleted {
		t.Fatalf("expected cache.DeleteByID to be called")
	}
}

func TestUserService_Delete(t *testing.T) {
	cacheDeleted := false

	repo := &mock.MockUserRepository{
		GetByIDFn: func(id string) (*models.User, error) {
			return &models.User{Base: models.Base{ID: uuid.New()}}, nil
		},
		DeleteFn: func(id string) error {
			return nil
		},
	}

	cache := &mock.MockUserCache{
		DeleteByIDFn: func(ctx context.Context, id string) error {
			cacheDeleted = true
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	svc := NewUserService(repo, cache, logger)
	err := svc.Delete("1")

	if err != nil {
		t.Fatalf("unexpected error")
	}
	if !cacheDeleted {
		t.Fatalf("expected cache.DeleteByID")
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old"), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		user    *models.User
		oldPass string
		newPass string
		wantErr string
	}{
		{
			name:    "user not found",
			wantErr: "User not found",
		},
		{
			name:    "wrong old password",
			user:    &models.User{Base: models.Base{ID: uuid.New()}, Password: string(oldHash)},
			oldPass: "wrong",
			newPass: "new",
			wantErr: "старый пароль неверен",
		},
		{
			name:    "empty new password",
			user:    &models.User{Base: models.Base{ID: uuid.New()}, Password: string(oldHash)},
			oldPass: "old",
			newPass: " ",
			wantErr: "новый пароль не должен быть пустым",
		},
		{
			name:    "success",
			user:    &models.User{Base: models.Base{ID: uuid.New()}, Password: string(oldHash)},
			oldPass: "old",
			newPass: "new",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDeleted := false

			repo := &mock.MockUserRepository{
				GetByIDFn: func(id string) (*models.User, error) {
					if tt.user == nil {
						return nil, gorm.ErrRecordNotFound
					}
					return tt.user, nil
				},
				UpdateFn: func(u *models.User) error {
					return nil
				},
			}

			cacheMock := &mock.MockUserCache{
				DeleteByIDFn: func(ctx context.Context, id string) error {
					cacheDeleted = true
					return nil
				},
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
			svc := NewUserService(repo, cacheMock, logger)
			err := svc.ChangePassword("1", tt.oldPass, tt.newPass)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !cacheDeleted {
				t.Fatalf("expected cache.DeleteByID")
			}
		})
	}
}
