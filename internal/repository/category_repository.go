package repository

import (
	"context"
	"effective-project/internal/models"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error

	List(ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID) ([]models.Category, error)

	GetByID(id string) (*models.Category, error)

	Update(category *models.Category) error

	Delete(id string) error
}

type gormCategoryRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewCategoryRepository(db *gorm.DB, logger *slog.Logger) CategoryRepository {
	return &gormCategoryRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormCategoryRepository) Create(category *models.Category) error {
	op := "repository.category.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("category", category),
	)

	if err := r.DB.Create(category).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormCategoryRepository) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]models.Category, error) {
	op := "repository.category.get_all"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	rows := make([]models.Category, 0, limit)

	q := r.DB.
		WithContext(ctx).
		Table("category").
		Select(`
	id,
	name,
	created_at
	`).
		Order("created_at ASC").
		Order("id ASC").
		Limit(limit)

	if lastCreatedAt != nil && lastID != nil {
		q = q.Where(`
			(created_at > ?)
			OR (created_at = ? AND id > ?)
		`, *lastCreatedAt, *lastCreatedAt, *lastID)
	}

	if lastCreatedAt != nil && lastID != nil {
		q = q.Where(`
			(created_at > ?)
			OR (created_at = ? AND id > ?)
		`, *lastCreatedAt, *lastCreatedAt, *lastID)
	}

	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *gormCategoryRepository) GetByID(id string) (*models.Category, error) {
	op := "repository.category.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	var category models.Category
	if err := r.DB.First(&category, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return &category, nil
}

func (r *gormCategoryRepository) Update(category *models.Category) error {
	op := "repository.category.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("category", category),
	)

	if err := r.DB.Save(category).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormCategoryRepository) Delete(id string) error {
	op := "repository.category.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	if err := r.DB.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}
