package service

import (
	"context"
	cache "effective-project/internal/cache"
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"effective-project/internal/repository"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req *dto.UserCreateRequest) (*models.User, error)

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID) ([]models.User, error)

	GetByID(id string) (*models.User, error)

	Update(id string, user *dto.UserUpdateRequest) (*models.User, error)

	Delete(id string) error

	ChangePassword(userID string, oldPassword, newPassword string) error
}

type userService struct {
	repo   repository.UserRepository
	cache  cache.UserCache
	logger *slog.Logger
}

func NewUserService(
	repo repository.UserRepository,
	cache cache.UserCache,
	logger *slog.Logger,
) UserService {
	return &userService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *userService) Create(req *dto.UserCreateRequest) (*models.User, error) {
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if err := s.repo.Create(user); err != nil {
		s.logger.Error("service.user.create: failed to create user:", slog.Any("error", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID) ([]models.User, error) {
	users, err := s.repo.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		s.logger.Error("service.user.get_all: failed to get users:", slog.Any("error", err))
		return nil, err
	}

	return users, nil
}

func (s *userService) GetByID(id string) (*models.User, error) {
	ctx := context.Background()

	if user, err := s.cache.GetByID(ctx, id); err == nil {
		s.logger.Debug("user cache hit", "user_id", id)
		return user, nil
	}

	s.logger.Debug("user cache miss", "user_id", id)

	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("service.user.get_by_id: failed to get user:", slog.Any("error", err))
		return nil, err
	}

	if err := s.cache.Set(ctx, user, time.Minute*10); err != nil {
		s.logger.Warn("service.user.get_by_id: failed to set cache", slog.Any("error", err))
	}

	return user, nil
}

func (s *userService) Update(id string, req *dto.UserUpdateRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("service.user.update: failed to get user:", slog.Any("error", err))
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repo.Update(user); err != nil {
		s.logger.Error("service.user.update: failed to update user:", slog.Any("error", err))
		return nil, err
	}

	userID := user.ID.String()

	ctx := context.Background()

	if err := s.cache.DeleteByID(ctx, userID); err != nil {
		s.logger.Warn("service.user.update: failed to delete cache by id", slog.Any("error", err))
	}

	return user, nil
}

func (s *userService) Delete(id string) error {
	user, _ := s.repo.GetByID(id)

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("service.user.delete: failed to delete user:", slog.Any("error", err))
		return err
	}

	if user != nil {
		userID := user.ID.String()
		ctx := context.Background()

		if err := s.cache.DeleteByID(ctx, userID); err != nil {
			s.logger.Warn("service.user.delete: failed to delete cache by id", slog.Any("error", err))
		}
	}

	return nil
}

func checkPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func (s *userService) ChangePassword(
	userID string,
	oldPassword,
	newPassword string,
) error {
	s.logger.Debug("ChangePassword called", "user_id", userID)
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found for ChangePassword", "user_id", userID)
			return errors.New("User not found")
		}
		s.logger.Error("error fetching user for ChangePassword", "error", err, "user_id", userID)
		return err
	}

	if err := checkPassword(user.Password, oldPassword); err != nil {
		s.logger.Warn("old password mismatch", "user_id", userID)
		return errors.New("старый пароль неверен")
	}

	if strings.TrimSpace(newPassword) == "" {
		s.logger.Warn("new password is empty", "user_id", userID)
		return errors.New("новый пароль не должен быть пустым")
	}

	hashed, err := hashPassword(newPassword)
	if err != nil {
		s.logger.Error("failed to hash new password", "error", err, "user_id", userID)
		return err
	}

	user.Password = hashed

	if err := s.repo.Update(user); err != nil {
		s.logger.Error("failed to update user password in repo", "error", err, "user_id", userID)
		return err
	}

	id := user.ID.String()

	ctx := context.Background()

	_ = s.cache.DeleteByID(ctx, id)
	s.logger.Info("password changed", "user_id", userID)
	return nil
}

func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 14)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
