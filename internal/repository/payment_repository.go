package repository

import (
	"context"
	"time"

	"effective-project/internal/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(service *models.Payment) error

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID,
	) ([]models.Payment, error)

	GetByID(id string) (*models.Payment, error)

	Update(service *models.Payment) error

	Delete(id string) error
}

type gormPaymentRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewPaymentRepository(db *gorm.DB, logger *slog.Logger) PaymentRepository {
	return &gormPaymentRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormPaymentRepository) Create(payment *models.Payment) error {
	op := "repository.payment.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("payment", payment),
	)

	if err := r.DB.Create(payment).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormPaymentRepository) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]models.Payment, error) {

	op := "repository.payment.list"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Int("limit", limit),
		slog.Any("lastCreatedAt", lastCreatedAt),
		slog.Any("lastID", lastID),
	)

	payments := make([]models.Payment, 0, limit)

	q := r.DB.
		WithContext(ctx).
		Model(&models.Payment{}).
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

	if err := q.Find(&payments).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return payments, nil
}

func (r *gormPaymentRepository) GetByID(id string) (*models.Payment, error) {
	op := "repository.payment.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("id", id),
	)

	var payment models.Payment
	if err := r.DB.First(&payment, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return &payment, nil
}

func (r *gormPaymentRepository) Update(payment *models.Payment) error {
	op := "repository.payment.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("payment", payment),
	)

	if err := r.DB.Save(payment).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormPaymentRepository) Delete(id string) error {
	op := "repository.payment.delete"

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
