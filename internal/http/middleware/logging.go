package middleware

import (
	"effective-project/internal/dto"
	"effective-project/internal/models"
	"effective-project/internal/service"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	auth   service.AuthService
	users  service.UserService
	logger *slog.Logger
}

func NewAuthHandler(auth service.AuthService, users service.UserService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		auth:   auth,
		users:  users,
		logger: logger,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, jwtCfg service.JWTConfig) {
	auth := r.Group("/auth")
	// ----публичные----
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// -----нужна-авторизация----
	protected := auth.Group("")
	protected.Use(AuthMiddleware(jwtCfg))
	protected.POST("/refresh", h.Refresh)
	protected.POST("/logout", h.Logout)
	protected.GET("/me", h.Me)
	protected.PUT("/me", h.UpdateMe)
	protected.PUT("/me/password", h.ChangePassword)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Auth.Register", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	req.Role = models.RoleUser

	user, err := h.users.Create(&req)

	if err != nil {
		h.logger.Error("Ошибка создания пользователя через /register", "error", err.Error(), "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Регистрация пользователя выполнена", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Auth.Login", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	token, err := h.auth.Login(req.Email, req.Password)

	if err != nil {
		if err == service.ErrInvalidCredentials {
			h.logger.Warn("Неудачная попытка входа", "email", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный email или пароль",
			})
			return
		}
		h.logger.Error("Ошибка при попытке логина", "error", err.Error(), "email", req.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Пользователь вошёл", "email", req.Email)
	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}
	user, err := h.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, errors.New("user not found")) {
			h.logger.Warn("Пользователь не найден в Auth.Me", "user_id", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка получения профиля в Auth.Me", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Профиль получен", "user_id", user.ID)
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}
	var req dto.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Auth.UpdateMe", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	user, err := h.users.Update(userID, &req)

	if err != nil {
		if errors.Is(err, errors.New("user not found")) {
			h.logger.Warn("Пользователь не найден в Auth.UpdateMe", "user_id", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка обновления профиля в Auth.UpdateMe", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Профиль обновлён", "user_id", user.ID)
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	idVal, exists := c.Get("userID")
	roleVal, existsRole := c.Get("userRole")

	if !exists || !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userdID, okID := idVal.(uuid.UUID)
	role, okRole := roleVal.(string)

	if !okID || !okRole {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный данные",
		})
		return
	}

	token, err := h.auth.GenerateToken(userdID, role)

	if err != nil {
		h.logger.Error("Ошибка генерации токена в Auth.Refresh", "error", err.Error(), "user_id", userdID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Токен обновлён", "user_id", userdID)
	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}

	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Auth.ChangePassword", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	if err := h.users.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		h.logger.Error("Ошибка смены пароля", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Пароль успешно изменён", "user_id", userID)
	c.Status(http.StatusOK)
}
