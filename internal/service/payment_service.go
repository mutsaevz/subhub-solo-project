package service

import (
	"context"
	"effective-project/internal/cache"
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"effective-project/internal/repository"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	Create(req *dto.PaymentCreateRequest) (models.Payment, error)

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID,
	) ([]models.Payment, error)

	GetByID(id string) (*models.Payment, error)

	Update(id string, req *dto.PaymentUpdateRequest) (*models.Payment, error)

	Delete(id string) error
}

type paymentService struct {
	paymentRepo  repository.PaymentRepository
	paymentCache cache.PaymentCache
	logger       *slog.Logger
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	paymentCache cache.PaymentCache,
	logger *slog.Logger,
) PaymentService {
	return &paymentService{
		paymentRepo:  paymentRepo,
		paymentCache: paymentCache,
		logger:       logger,
	}
}

func (s *paymentService) Create(req *dto.PaymentCreateRequest) (models.Payment, error) {
	payment := models.Payment{
		SubscriptionID: req.SubscriptionID,
		OrderID:        req.OrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		PaidAt:         req.PaidAt,
		PaymentStatus:  req.PaymentStatus,
		Provider:       req.Provider,
	}

	if err := s.paymentRepo.Create(&payment); err != nil {
		s.logger.Error("service.payment.create: failed to create payment", slog.Any("error", err))
		return models.Payment{}, err
	}

	return payment, nil
}

func (s *paymentService) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID,
) ([]models.Payment, error) {
	payments, err := s.paymentRepo.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		s.logger.Error("service.payment.get_all: failed to get payments", slog.Any("error", err))
		return nil, err
	}

	return payments, nil
}

func (s *paymentService) GetByID(id string) (*models.Payment, error) {
	ctx := context.Background()

	payment, err := s.paymentCache.GetByID(ctx, id)
	if err == nil {
		return payment, nil
	}

	payment, err = s.paymentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.payment.get_by_id: failed to get payment", slog.Any("error", err))
		return nil, err
	}

	if err := s.paymentCache.Set(ctx, payment, time.Minute*10); err != nil {
		s.logger.Warn("service.payment.get_by_id: failed to set cache", slog.Any("error", err))
	}

	return payment, nil
}

func (s *paymentService) Update(id string, req *dto.PaymentUpdateRequest) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.payment.update: failed to get payment", slog.Any("error", err))
		return nil, err
	}

	if req.Amount != nil {
		payment.Amount = *req.Amount
	}
	if req.Currency != nil {
		payment.Currency = *req.Currency
	}
	if req.PaidAt != nil {
		payment.PaidAt = *req.PaidAt
	}
	if req.PaymentStatus != nil {
		payment.PaymentStatus = *req.PaymentStatus
	}
	if req.Provider != nil {
		payment.Provider = *req.Provider
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		s.logger.Error("service.payment.update: failed to update payment", slog.Any("error", err))
		return nil, err
	}

	if err := s.paymentCache.Set(context.Background(), payment, time.Minute*10); err != nil {
		s.logger.Warn("service.payment.update: failed to update cache", slog.Any("error", err))
	}

	return payment, nil
}

func (s *paymentService) Delete(id string) error {
	if err := s.paymentRepo.Delete(id); err != nil {
		s.logger.Error("service.payment.delete: failed to delete payment", slog.Any("error", err))
		return err
	}

	if err := s.paymentCache.Delete(context.Background(), id); err != nil {
		s.logger.Warn("service.payment.delete: failed to delete cache", slog.Any("error", err))
	}

	return nil
}
