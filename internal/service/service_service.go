package service

import (
	"context"
	"effective-project/internal/cache"
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"effective-project/internal/repository"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ServiceService interface {
	Create(req *dto.ServiceCreateRequest) (*models.Service, error)

	List(
		ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID) ([]models.Service, error)

	GetByID(id string) (*models.Service, error)

	Update(id string, service *dto.ServiceUpdateRequest) (*models.Service, error)

	Delete(id string) error
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
	cache       cache.Cache
	logger      *slog.Logger
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	cache cache.Cache,
	logger *slog.Logger,
) ServiceService {
	return &serviceService{
		serviceRepo: serviceRepo,
		cache:       cache,
		logger:      logger,
	}
}

func (s *serviceService) Create(req *dto.ServiceCreateRequest) (*models.Service, error) {
	service := &models.Service{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		LogoUrl:    req.LogoUrl,
		Website:    req.Website,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		s.logger.Error("service.service.create: failed to create service", slog.Any("error", err))
		return nil, err
	}

	_ = s.cache.Set(context.Background(), fmt.Sprintf("service:%s", service.ID.String()), service)

	return service, nil
}

func (s *serviceService) List(
	ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID) ([]models.Service, error) {
	services, err := s.serviceRepo.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		s.logger.Error("service.service.get_all: failed to get services", slog.Any("error", err))
		return nil, err
	}

	return services, nil
}

func (s *serviceService) GetByID(id string) (*models.Service, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("service:%s", id)

	var cached models.Service
	if ok, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && ok {
		return &cached, nil
	}

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.service.get_by_id: failed to get service", slog.Any("error", err))
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, service)
	return service, nil
}

func (s *serviceService) Update(id string, req *dto.ServiceUpdateRequest) (*models.Service, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.service.update: failed to get service", slog.Any("error", err))
		return nil, err
	}

	if req.Name != nil {
		service.Name = *req.Name
	}
	if req.LogoUrl != nil {
		service.LogoUrl = *req.LogoUrl
	}
	if req.Website != nil {
		service.Website = *req.Website
	}

	if err := s.serviceRepo.Update(service); err != nil {
		s.logger.Error("service.service.update: failed to update service", slog.Any("error", err))
		return nil, err
	}

	_ = s.cache.Delete(context.Background(), fmt.Sprintf("service:%s", id))

	return service, nil
}

func (s *serviceService) Delete(id string) error {
	if err := s.serviceRepo.Delete(id); err != nil {
		s.logger.Error("service.service.delete: failed to delete service", slog.Any("error", err))
		return err
	}

	_ = s.cache.Delete(context.Background(), fmt.Sprintf("service:%s", id))

	return nil
}
