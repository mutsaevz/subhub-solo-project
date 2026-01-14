package repository

import (
	"effective-project/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error

	GetByID(id string) (*models.Order, error)

	Update(order *models.Order) error
}

type gormOrderRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewOrderRepository(db *gorm.DB, logger *slog.Logger) OrderRepository {
	return &gormOrderRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormOrderRepository) Create(order *models.Order) error {
	op := "repository.order.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("order", order),
	)

	if err := r.db.Create(order).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormOrderRepository) GetByID(id string) (*models.Order, error) {
	op := "repository.order.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("id", id),
	)

	var order models.Order
	if err := r.db.First(&order, "id = ?", id).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return &order, nil
}

func (r *gormOrderRepository) Update(order *models.Order) error {
	op := "repository.order.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Any("order", order),
	)

	if err := r.db.Save(order).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}
