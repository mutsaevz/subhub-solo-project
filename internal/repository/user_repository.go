package repository

import (
	"context"
	"effective-project/internal/models"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID,
	) ([]models.User, error)

	GetByID(id string) (*models.User, error)

	GetByEmail(email string) (*models.User, error)

	Update(user *models.User) error

	Delete(id string) error
}

type gormUserRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &gormUserRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormUserRepository) Create(user *models.User) error {
	op := "repository.user.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("user", user),
	)

	if err := r.DB.Create(user).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormUserRepository) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]models.User, error) {
	op := "repository.user.get_all"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	user := make([]models.User, 0, limit)

	q := r.DB.
		WithContext(ctx).
		Model(&models.User{}).
		Order("created_at ASC").
		Order("id ASC").
		Limit(limit)

	if lastCreatedAt != nil && lastID != nil {
		q = q.Where(
			"(created_at > ?) OR (created_at = ? AND id > ?)",
			*lastCreatedAt,
			*lastCreatedAt,
			*lastID,
		)
	}

	if err := q.Find(&user).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return user, nil
}

func (r *gormUserRepository) GetByID(id string) (*models.User, error) {
	op := "repository.user.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	var user models.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	op := "repository.user.get_by_email"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("email", email),
	)

	var user models.User
	if err := r.DB.First(&user, "email = ?", email).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) Update(user *models.User) error {
	op := "repository.user.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("user", user),
	)

	if err := r.DB.Save(user).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormUserRepository) Delete(id string) error {
	op := "repository.user.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	if err := r.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}
