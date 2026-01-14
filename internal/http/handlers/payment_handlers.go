package handlers

import (
	"effective-project/internal/dto"
	"effective-project/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandlers struct {
	paymentService service.PaymentService
	logger         *slog.Logger
}

func NewPaymentHandlers(paymentService service.PaymentService, logger *slog.Logger) *PaymentHandlers {
	return &PaymentHandlers{
		paymentService: paymentService,
		logger:         logger,
	}
}

// RegisterRoutes регистрирует роуты в Gin
func (h *PaymentHandlers) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/payments", h.Create)
	r.GET("/payments", h.List)
	r.GET("/payments/:id", h.GetByID)
	r.PUT("/payments/:id", h.Update)
	r.DELETE("/payments/:id", h.Delete)
}

func (h *PaymentHandlers) Create(c *gin.Context) {
	var req dto.PaymentCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("payment.create: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	payment, err := h.paymentService.Create(&req)
	if err != nil {
		h.logger.Error("payment.create: failed to create payment", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandlers) List(c *gin.Context) {
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

	payments, err := h.paymentService.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		h.logger.Error("payment.list: failed to list payments", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments"})
		return
	}

	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
	}

	var nextCursor *cursor
	if len(payments) == limit {
		last := payments[len(payments)-1]
		nextCursor = &cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       payments,
		"next_cursor": nextCursor,
	})
}

func (h *PaymentHandlers) GetByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := h.paymentService.GetByID(id)
	if err != nil {
		h.logger.Error("payment.getByID: failed to get payment", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandlers) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.PaymentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("payment.update: failed to decode request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	payment, err := h.paymentService.Update(id, &req)
	if err != nil {
		h.logger.Error("payment.update: failed to update payment", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandlers) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.paymentService.Delete(id); err != nil {
		h.logger.Error("payment.delete: failed to delete payment", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
		return
	}

	c.Status(http.StatusNoContent)
}
