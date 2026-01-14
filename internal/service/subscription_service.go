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

type SubscriptionService interface {
	Create(req *dto.SubscriptionCreateRequest) (*models.Subscription, error)

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID) ([]dto.SubscriptionResponse, error)

	GetByID(id string) (*dto.SubscriptionResponse, error)

	Update(id string, req *dto.SubscriptionUpdateRequest) (*models.Subscription, error)

	Delete(id string) error

	CalculateTotal(
		ctx context.Context,
		f dto.TotalFilter,
	) (int, error)
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
	serviceRepo      repository.ServiceRepository
	paymentRepo      repository.PaymentRepository

	subscriptionCache cache.SubscriptionCache
	logger            *slog.Logger
}

func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	serviceRepo repository.ServiceRepository,
	paymentRepo repository.PaymentRepository,
	subscriptionCache cache.SubscriptionCache,
	logger *slog.Logger,
) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo:  subscriptionRepo,
		serviceRepo:       serviceRepo,
		paymentRepo:       paymentRepo,
		subscriptionCache: subscriptionCache,
		logger:            logger,
	}
}

func (s *subscriptionService) Create(req *dto.SubscriptionCreateRequest) (*models.Subscription, error) {

	var subscription = &models.Subscription{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Price:     req.Price,
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		s.logger.Error("service.subscription.create: failed to create subscription", slog.Any("error", err))
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID) ([]dto.SubscriptionResponse, error) {

	subscriptions, err := s.subscriptionRepo.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		s.logger.Error("service.subscription.get_all: failed to get subscriptions", slog.Any("error", err))
		return nil, err
	}

	return subscriptions, nil
}

func (s *subscriptionService) GetByID(id string) (*dto.SubscriptionResponse, error) {
	ctx := context.Background()

	// 1. Пытаемся взять из кеша (модель)
	if sub, err := s.subscriptionCache.GetByID(ctx, id); err == nil {
		return &dto.SubscriptionResponse{
			ID:        sub.ID,
			ServiceID: sub.ServiceID,
			StartDate: sub.StartDate,
			EndDate:   sub.EndDate,
			Price:     sub.Price,
			CreatedAt: sub.CreatedAt,
		}, nil
	}

	// 2. Берём из репозитория (DTO)
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		s.logger.Error(
			"service.subscription.get_by_id: failed to get subscription",
			slog.Any("error", err),
		)
		return nil, err
	}

	// 3. Кладём в кеш модель
	model, err := s.subscriptionRepo.GetModelByID(id)
	if err == nil {
		if err := s.subscriptionCache.Set(ctx, model, time.Minute*10); err != nil {
			s.logger.Warn(
				"service.subscription.get_by_id: failed to set cache",
				slog.Any("error", err),
			)
		}
	}

	return subscription, nil
}

func (s *subscriptionService) Update(id string, req *dto.SubscriptionUpdateRequest) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetModelByID(id)
	if err != nil {
		s.logger.Error("service.subscription.update: failed to get subscription", slog.Any("error", err))
		return nil, err
	}

	if req.EndDate != nil {
		subscription.EndDate = req.EndDate
	}
	if req.StartDate != nil {
		subscription.StartDate = *req.StartDate
	}
	if req.Price != nil {
		subscription.Price = *req.Price
	}

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("service.subscription.update: failed to update subscription", slog.Any("error", err))
		return nil, err
	}

	ctx := context.Background()
	_ = s.subscriptionCache.DeleteByID(ctx, id)

	return subscription, nil
}

func (s *subscriptionService) Delete(id string) error {
	if err := s.subscriptionRepo.Delete(id); err != nil {
		s.logger.Error("service.subscription.delete: failed to cancel subscription", slog.Any("error", err))
		return err
	}

	_ = s.subscriptionCache.DeleteByID(context.Background(), id)

	return nil
}

func (s *subscriptionService) CalculateTotal(
	ctx context.Context,
	f dto.TotalFilter,
) (int, error) {

	rows, err := s.subscriptionRepo.FindForTotal(ctx, f)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, row := range rows {
		months := countMonths(
			row.StartDate,
			row.EndDate,
			f.From,
			f.To,
		)
		total += months * row.Price
	}

	return total, nil
}

func countMonths(
	subStart time.Time,
	subEnd *time.Time,
	periodFrom time.Time,
	periodTo time.Time,
) int {

	start := maxDate(subStart, periodFrom)

	end := periodTo
	if subEnd != nil && subEnd.Before(end) {
		end = *subEnd
	}

	if start.After(end) {
		return 0
	}

	return (end.Year()-start.Year())*12 +
		int(end.Month()-start.Month()) + 1
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
