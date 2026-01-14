package handlers

import (
	"net/http"
	"strconv"
	"time"

	"effective-project/internal/dto"
	"effective-project/internal/http/middleware"
	"effective-project/internal/service"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")

	admin := users.Group("")
	admin.Use(middleware.RequireRole("admin"))

	admin.POST("", h.Create)
	admin.GET("/email", h.GetByEmail)
	admin.GET("", h.List)
	admin.GET("/:id", h.GetByID)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.user.create: invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Create(&req)
	if err != nil {
		h.logger.Error("handler.user.create: failed to create user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) List(c *gin.Context) {
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

	users, err := h.userService.List(ctx, limit, lastCreatedAt, lastID)
	if err != nil {
		h.logger.Error("handler.user.list: failed to list users", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
	}

	var nextCursor *cursor
	if len(users) == limit {
		last := users[len(users)-1]
		nextCursor = &cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       users,
		"next_cursor": nextCursor,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
	if err != nil {
		h.logger.Error("handler.user.get_by_id: failed to get user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler.user.update: invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Update(id, &req)
	if err != nil {
		h.logger.Error("handler.user.update: failed to update user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.Delete(id); err != nil {
		h.logger.Error("handler.user.delete: failed to delete user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Query("email")

	user, err := h.userService.GetByEmail(email)
	if err != nil {
		h.logger.Error("handler.user.get_by_email: failed to get user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
