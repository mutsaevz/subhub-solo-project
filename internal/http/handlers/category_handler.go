package handlers

import (
	"effective-project/internal/dto"
	"effective-project/internal/http/middleware"
	"effective-project/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	logger          *slog.Logger
}

func NewCategoryHandler(categoryService service.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.RouterGroup) {
	categories := r.Group("/categories")
	// Public routes
	categories.GET("", h.List)
	categories.GET("/:id", h.GetByID)

	admin := categories.Group("")
	admin.Use(middleware.RequireRole("admin"))

	// Admin routes
	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handlers.category.create: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category, err := h.categoryService.Create(&req)
	if err != nil {
		h.logger.Error("handlers.category.create: failed to create category", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	var (
		lastCreatedAt *time.Time
		lastID        *uuid.UUID
	)

	if v := c.Query("created_at"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid created_at"})
			return
		}
		lastCreatedAt = &t
	}

	if v := c.Query("id"); v != "" {
		uid, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		lastID = &uid
	}

	if (lastCreatedAt == nil) != (lastID == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "created_at and id must be used together"})
		return
	}

	categories, err := h.categoryService.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		h.logger.Error("handlers.category.list: failed to list categories", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}

	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
	}

	var nextCursor *cursor
	if len(categories) == limit {
		last := categories[len(categories)-1]
		nextCursor = &cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       categories,
		"next_cursor": nextCursor,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		h.logger.Error("handlers.category.getByID: failed to get category", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handlers.category.update: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category, err := h.categoryService.Update(id, &req)
	if err != nil {
		h.logger.Error("handlers.category.update: failed to update category", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryService.Delete(id); err != nil {
		h.logger.Error("handlers.category.delete: failed to delete category", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.Status(http.StatusNoContent)
}
