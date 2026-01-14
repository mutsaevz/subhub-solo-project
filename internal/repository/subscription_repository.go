package repository

import (
	"context"
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID,
	) ([]dto.SubscriptionResponse, error)

	GetByID(id string) (*dto.SubscriptionResponse, error)

	Update(subscription *models.Subscription) error

	Delete(id string) error

	FindForTotal(
		ctx context.Context,
		f dto.TotalFilter,
	) ([]dto.SubscriptionRow, error)

	GetModelByID(id string) (*models.Subscription, error)
}

type gormSubscriptionRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewSubscriptionRepository(db *gorm.DB, logger *slog.Logger) SubscriptionRepository {
	return &gormSubscriptionRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormSubscriptionRepository) Create(subscription *models.Subscription) error {
	op := "repository.subscription.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("subscription", subscription),
	)

	if err := r.DB.Create(subscription).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormSubscriptionRepository) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]dto.SubscriptionResponse, error) {

	op := "repository.subscription.get_all"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	result := make([]dto.SubscriptionResponse, 0, limit)

	q := r.DB.WithContext(ctx).
		Table("subscriptions").
		Select(`
			subscriptions.id,
			subscriptions.start_date,
			subscriptions.end_date,
			subscriptions.price,
			subscriptions.service_id,
			services.name AS service_name
		`).
		Joins("JOIN services ON services.id = subscriptions.service_id").
		Order("subscriptions.created_at ASC").
		Order("subscriptions.id ASC").
		Limit(limit)

	if lastCreatedAt != nil && lastID != nil {
		q = q.Where(
			"(subscriptions.created_at > ?) OR (subscriptions.created_at = ? AND subscriptions.id > ?)",
			*lastCreatedAt,
			*lastCreatedAt,
			*lastID,
		)
	}

	if err := q.Scan(&result).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return result, nil
}

func (r *gormSubscriptionRepository) GetByID(id string) (*dto.SubscriptionResponse, error) {
	op := "repository.subscription.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	var result dto.SubscriptionResponse

	if err := r.DB.
		Table("subscriptions").
		Select(`
			subscriptions.id,
			subscriptions.start_date,
			subscriptions.end_date,
			subscriptions.price,
			subscriptions.service_id,
			services.name AS service_name
		`).
		Joins("JOIN services ON services.id = subscriptions.service_id").
		Where("subscriptions.id = ?", id).
		Scan(&result).Error; err != nil {

		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &result, nil
}

func (r *gormSubscriptionRepository) Update(subscription *models.Subscription) error {
	op := "repository.subscription.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("subscription", subscription),
	)

	if err := r.DB.Save(subscription).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormSubscriptionRepository) Delete(id string) error {
	op := "repository.subscription.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	if err := r.DB.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormSubscriptionRepository) FindForTotal(
	ctx context.Context,
	f dto.TotalFilter,
) ([]dto.SubscriptionRow, error) {

	q := r.DB.WithContext(ctx).
		Table("subscriptions").
		Select(`
			subscriptions.start_date,
			subscriptions.end_date,
			subscriptions.price,
			services.name AS service_name
		`).
		Joins("JOIN services ON services.id = subscriptions.service_id").
		Where("subscriptions.start_date <= ?", f.To).
		Where("(subscriptions.end_date IS NULL OR subscriptions.end_date >= ?)", f.From)

	if f.UserID.String() != "" {
		q = q.Where("subscriptions.user_id = ?", f.UserID)
	}

	if f.ServiceName != "" {
		q = q.Where("services.name = ?", f.ServiceName)
	}

	var rows []dto.SubscriptionRow
	err := q.Scan(&rows).Error
	return rows, err
}

func (r *gormSubscriptionRepository) GetModelByID(id string) (*models.Subscription, error) {
	op := "repository.subscription.get_model_by_id"
	var subscription models.Subscription

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("id", id),
	)

	if err := r.DB.First(&subscription, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return &subscription, nil
}
