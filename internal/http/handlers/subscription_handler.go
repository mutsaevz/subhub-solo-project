package handlers

import (
	"effective-project/internal/dto"
	"effective-project/internal/http/middleware"
	"effective-project/internal/service"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	logger              *slog.Logger
}

func NewSubscriptionHandler(
	subscriptionService service.SubscriptionService,
	logger *slog.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(r *gin.RouterGroup) {
	subscriptions := r.Group("/subscriptions")

	admin := subscriptions.Group("")
	admin.Use(middleware.RequireRole("admin"))

	subscriptions.POST("", h.Create)
	subscriptions.GET("/total", h.GetTotal)
	subscriptions.GET("", h.List)
	subscriptions.GET("/:id", h.GetByID)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req dto.SubscriptionCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.subscription.create: invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	subscription, err := h.subscriptionService.Create(&req)
	if err != nil {
		h.logger.Error("handler.subscription.create: failed to create subscription", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	// limit
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "created_at and id must be used together",
		})
		return
	}

	subscriptions, err := h.subscriptionService.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		h.logger.Error("handler.subscription.list: failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
	}

	var nextCursor *cursor
	if len(subscriptions) == limit {
		last := subscriptions[len(subscriptions)-1]
		nextCursor = &cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       subscriptions,
		"next_cursor": nextCursor,
	})
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	subscription, err := h.subscriptionService.GetByID(id)
	if err != nil {
		h.logger.Error("handler.subscription.get_by_id: failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.SubscriptionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.subscription.update: invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	subscription, err := h.subscriptionService.Update(id, &req)
	if err != nil {
		h.logger.Error("handler.subscription.update: failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.subscriptionService.Delete(id); err != nil {
		h.logger.Error("handler.subscription.delete: failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) GetTotal(c *gin.Context) {
	from, err := parseMonth(c.Query("from"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid from"})
		return
	}

	to, err := parseMonth(c.Query("to"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid to"})
		return
	}

	var userID uuid.UUID
	if v := c.Query("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user_id"})
			return
		}
		userID = id
	}

	var serviceName string
	if v := c.Query("service_name"); v != "" {
		serviceName = v
	}

	total, err := h.subscriptionService.CalculateTotal(c.Request.Context(), dto.TotalFilter{
		From:        from,
		To:          to,
		UserID:      userID,
		ServiceName: serviceName,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"total_price": total})
}

func parseMonth(value string) (time.Time, error) {
	t, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		t.Year(),
		t.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	), nil
}
