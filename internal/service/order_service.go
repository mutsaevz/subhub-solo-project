package service

import (
	"context"
	"time"

	"effective-project/internal/cache"
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"effective-project/internal/repository"
	"log/slog"
)

type OrderService interface {
	Create(req dto.OrderCreateRequest) (*models.Order, error)

	GetByID(id string) (*models.Order, error)

	Update(id string, req dto.OrderUpdateRequest) (*models.Order, error)
}

type orderService struct {
	orderRepo  repository.OrderRepository
	orderCache cache.OrderCache
	logger     *slog.Logger
}

func NewOrderService(orderRepo repository.OrderRepository, orderCache cache.OrderCache, logger *slog.Logger) OrderService {
	return &orderService{
		orderRepo:  orderRepo,
		orderCache: orderCache,
		logger:     logger,
	}
}

func (s *orderService) Create(req dto.OrderCreateRequest) (*models.Order, error) {
	op := "service.order.create"
	ctx := context.Background()

	s.logger.Debug("service call", slog.String("op", op))

	order := &models.Order{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		IsPaid:    req.IsPaid,
	}

	if err := s.orderRepo.Create(order); err != nil {
		s.logger.Error("service.order.create: failed to create order", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	if err := s.orderCache.Set(ctx, order, time.Minute*5); err != nil {
		s.logger.Warn("service.order.create: failed to set order in cache", slog.String("op", op), slog.Any("error", err))
	}

	return order, nil
}

func (s *orderService) GetByID(id string) (*models.Order, error) {
	op := "service.order.get_by_id"
	ctx := context.Background()

	s.logger.Debug("service call", slog.String("op", op), slog.Any("id", id))

	order, err := s.orderCache.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn("service.order.get_by_id: failed to get order from cache", slog.String("op", op), slog.Any("error", err))
	}
	if order != nil {
		return order, nil
	}

	order, err = s.orderRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.order.get_by_id: failed to get order", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	if err := s.orderCache.Set(ctx, order, time.Minute*5); err != nil {
		s.logger.Warn("service.order.get_by_id: failed to set order in cache", slog.String("op", op), slog.Any("error", err))
	}

	return order, nil
}

func (s *orderService) Update(id string, req dto.OrderUpdateRequest) (*models.Order, error) {
	op := "service.order.update"
	ctx := context.Background()

	s.logger.Debug("service call", slog.String("op", op))

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.order.update: failed to get order", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	if req.IsPaid != nil {
		order.IsPaid = *req.IsPaid
	}

	if err := s.orderRepo.Update(order); err != nil {
		s.logger.Error("service.order.update: failed to update order", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	if err := s.orderCache.Set(ctx, order, time.Minute*5); err != nil {
		s.logger.Warn("service.order.update: failed to set order in cache", slog.String("op", op), slog.Any("error", err))
	}

	return order, nil
}
