package repository

import (
	"context"
	"time"

	"effective-project/internal/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) error

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID,
	) ([]models.Service, error)

	GetByID(id string) (*models.Service, error)

	Update(service *models.Service) error

	Delete(id string) error
}

type gormServiceRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewServiceRepository(db *gorm.DB, logger *slog.Logger) ServiceRepository {
	return &gormServiceRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormServiceRepository) Create(service *models.Service) error {
	op := "repository.service.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("service", service),
	)

	if err := r.DB.Create(service).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormServiceRepository) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]models.Service, error) {

	op := "repository.service.list"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Int("limit", limit),
		slog.Any("lastCreatedAt", lastCreatedAt),
		slog.Any("lastID", lastID),
	)

	services := make([]models.Service, 0, limit)

	q := r.DB.
		WithContext(ctx).
		Model(&models.Service{}).
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

	if err := q.Find(&services).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return services, nil
}

func (r *gormServiceRepository) GetByID(id string) (*models.Service, error) {
	op := "repository.service.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	var service models.Service
	if err := r.DB.First(&service, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return &service, nil
}

func (r *gormServiceRepository) Update(service *models.Service) error {
	op := "repository.service.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("service", service),
	)

	if err := r.DB.Save(service).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormServiceRepository) Delete(id string) error {
	op := "repository.service.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	if err := r.DB.Delete(&models.Service{}, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}
