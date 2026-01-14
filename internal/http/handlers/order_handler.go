
package handlers

import (
	"net/http"

	"effective-project/internal/dto"
	"effective-project/internal/service"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type OrderHandler struct {
	orderService service.OrderService
	logger       *slog.Logger
}

func NewOrderHandler(
	orderService service.OrderService,
	logger *slog.Logger,
) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

// RegisterRoutes регистрирует роуты в gin.Engine или gin.RouterGroup
func (h *OrderHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/orders", h.Create)
	r.GET("/orders/:id", h.GetByID)
	r.PUT("/orders/:id", h.Update)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handlers.order.create: failed to decode request", slog.Any("error", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := h.orderService.Create(req)
	if err != nil {
		h.logger.Error("handlers.order.create: failed to create order", slog.Any("error", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderService.GetByID(id)
	if err != nil {
		h.logger.Error("handlers.order.get_by_id: failed to get order", slog.Any("error", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.OrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handlers.order.update: failed to decode request", slog.Any("error", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := h.orderService.Update(id, req)
	if err != nil {
		h.logger.Error("handlers.order.update: failed to update order", slog.Any("error", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, order)
}
