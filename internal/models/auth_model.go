package models

// Модель аутентификации пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"StrongPass123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8" example:"OldPass123"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"NewStrongPass123"`
}
