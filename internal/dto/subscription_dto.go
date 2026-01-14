package dto

import (
	"time"

	"github.com/google/uuid"
)

// DTO for filter
type TotalFilter struct {
	From        time.Time
	To          time.Time
	UserID      uuid.UUID
	ServiceName string
}

type SubFilter struct {
	From        time.Time
	To          time.Time
	UserID      string
	ServiceName string
}

// DTO для подписки пользователя на сервис
type SubscriptionCreateRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`

	ServiceID uuid.UUID `json:"service_id" binding:"required"`

	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   *time.Time `json:"end_date"`

	Price int `json:"price" binding:"required,gt=0" gorm:"not null;index"`
}

type SubscriptionUpdateRequest struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	Price *int `json:"price" binding:"required,gt=0" gorm:"not null;index"`
}

type SubscriptionRow struct {
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Price     int        `json:"price"`

	ServiceName string `json:"service_name"`
}

type SubscriptionResponse struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Price       int        `json:"price"`
	ServiceID   uuid.UUID  `json:"service_id"`
	ServiceName string     `json:"service_name"`
}
