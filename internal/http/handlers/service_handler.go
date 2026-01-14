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

type ServiceHandler struct {
	serviceService service.ServiceService
	logger         *slog.Logger
}

func NewServiceHandler(serviceService service.ServiceService, logger *slog.Logger) *ServiceHandler {
	return &ServiceHandler{
		serviceService: serviceService,
		logger:         logger,
	}
}

// RegisterRoutes registers service routes in Gin router
func (h *ServiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	services := r.Group("/services")
	// Public routes
	services.GET("", h.List)
	services.GET("/:id", h.GetByID)

	admin := services.Group("")
	admin.Use(middleware.RequireRole("admin"))

	// Admin routes
	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.ServiceCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.service.create: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	service, err := h.serviceService.Create(&req)
	if err != nil {
		h.logger.Error("handler.service.create: failed to create service", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	// limit
	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// cursor
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

	services, err := h.serviceService.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		h.logger.Error("handler.service.list: failed to list services", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list services"})
		return
	}

	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
	}

	var nextCursor *cursor
	if len(services) == limit {
		last := services[len(services)-1]
		nextCursor = &cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       services,
		"next_cursor": nextCursor,
	})
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	service, err := h.serviceService.GetByID(id)
	if err != nil {
		h.logger.Error("handler.service.get_by_id: failed to get service", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get service"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.ServiceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.service.update: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	service, err := h.serviceService.Update(id, &req)
	if err != nil {
		h.logger.Error("handler.service.update: failed to update service", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.serviceService.Delete(id); err != nil {
		h.logger.Error("handler.service.delete: failed to delete service", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.Status(http.StatusNoContent)
}
