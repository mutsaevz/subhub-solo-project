package dto

import (
	"effective-project/internal/models"
	"time"

	"github.com/google/uuid"
)

// DTO для платежа за подписку
type PaymentCreateRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" binding:"required"`
	OrderID        uuid.UUID `json:"order_id" binding:"required"`

	Amount   int    `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,len=3"`

	PaidAt time.Time `json:"paid_at"`

	PaymentStatus models.PaymentStatus `json:"payment_status" binding:"required"`
	Provider      string               `json:"provider" binding:"required,min=2,max=50"`
}

type PaymentUpdateRequest struct {
	Amount        *int                  `json:"amount" binding:"omitempty,gt=0"`
	Currency      *string               `json:"currency" binding:"omitempty,len=3"`
	PaidAt        *time.Time            `json:"paid_at"`
	PaymentStatus *models.PaymentStatus `json:"payment_status" binding:"omitempty"`
	Provider      *string               `json:"provider" binding:"omitempty,min=2,max=50"`
}
