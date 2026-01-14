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

type CategoryService interface {
	Create(category *dto.CategoryCreateRequest) (*models.Category, error)

	List(ctx context.Context,
		limit int,
		lastCreatedAt *time.Time,
		lastID *uuid.UUID) ([]models.Category, error)

	GetByID(id string) (*models.Category, error)

	Update(id string, category *dto.CategoryUpdateRequest) (*models.Category, error)

	Delete(id string) error
}

type categoryService struct {
	categoryRepo  repository.CategoryRepository
	categoryCache cache.CategoryCache
	logger       *slog.Logger
}

func NewCategoryService(categoryRepo repository.CategoryRepository, categoryCache cache.CategoryCache, logger *slog.Logger) CategoryService {
	return &categoryService{
		categoryRepo:  categoryRepo,
		categoryCache: categoryCache,
		logger:       logger,
	}
}

func (s *categoryService) Create(req *dto.CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("service.category.create: failed to create category", slog.Any("error", err))
		return nil, err
	}

	if err := s.categoryCache.Set(context.Background(), category, time.Minute*5); err != nil {
		s.logger.Warn("service.category.create: failed to set category cache", slog.Any("error", err))
	}

	return category, nil
}

func (s *categoryService) List(ctx context.Context,
	limit int,
	lastCreatedAt *time.Time,
	lastID *uuid.UUID) ([]models.Category, error) {
	categories, err := s.categoryRepo.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		s.logger.Error("service.category.get_all: failed to get categories", slog.Any("error", err))
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) GetByID(id string) (*models.Category, error) {
	category, err := s.categoryCache.GetByID(context.Background(), id)
	if err != nil {
		s.logger.Warn("service.category.get_by_id: failed to get category from cache", slog.Any("error", err))
	}
	if category != nil {
		return category, nil
	}

	category, err = s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.category.get_by_id: failed to get category", slog.Any("error", err))
		return nil, err
	}

	if err := s.categoryCache.Set(context.Background(), category, time.Minute*5); err != nil {
		s.logger.Warn("service.category.get_by_id: failed to set category cache", slog.Any("error", err))
	}

	return category, nil
}

func (s *categoryService) Update(id string, req *dto.CategoryUpdateRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Error("service.category.update: failed to get category", slog.Any("error", err))
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}

	if err := s.categoryRepo.Update(category); err != nil {
		s.logger.Error("service.category.update: failed to update category", slog.Any("error", err))
		return nil, err
	}

	if err := s.categoryCache.Set(context.Background(), category, time.Minute*5); err != nil {
		s.logger.Warn("service.category.update: failed to update category cache", slog.Any("error", err))
	}

	return category, nil
}

func (s *categoryService) Delete(id string) error {
	if err := s.categoryRepo.Delete(id); err != nil {
		s.logger.Error("service.category.delete: failed to delete category", slog.Any("error", err))
		return err
	}

	if err := s.categoryCache.Delete(context.Background(), id); err != nil {
		s.logger.Warn("service.category.delete: failed to delete category cache", slog.Any("error", err))
	}

	return nil
}
