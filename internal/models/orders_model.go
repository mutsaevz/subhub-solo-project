package models

import "github.com/google/uuid"

// Заказ, оформленный пользователем
type Order struct {
	Base

	UserID    uuid.UUID `json:"user_id" binding:"required" gorm:"type:uuid;not null;index"`
	ServiceID uuid.UUID `json:"service_id" binding:"required" gorm:"type:uuid;not null;index"`

	IsPaid bool `json:"is_paid" gorm:"not null;default:false;index"`
}
