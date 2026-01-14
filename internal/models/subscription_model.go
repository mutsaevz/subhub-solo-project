package models

import (
	"time"

	"github.com/google/uuid"
)

// Подписка пользователя на сервис
type Subscription struct {
	Base

	UserID uuid.UUID `json:"user_id" binding:"required" gorm:"type:uuid;not null;index"`
	User   User      `json:"-"`

	ServiceID uuid.UUID `json:"service_id" binding:"required" gorm:"type:uuid;not null;index"`
	Service   Service   `json:"-"`

	StartDate time.Time `json:"start_date" binding:"required" gorm:"not null;index"`
	EndDate   *time.Time `json:"end_date" gorm:"index"`

	Price    int    `json:"price" binding:"required,gt=0" gorm:"not null;index"`
}
