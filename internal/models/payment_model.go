package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentSucces PaymentStatus = "success"
	PaymentFailed PaymentStatus = "failed"
)

// Платеж за подписку
type Payment struct {
	Base

	SubscriptionID uuid.UUID    `json:"subscription_id" binding:"required" gorm:"type:uuid;not null;index"`
	Subscription   Subscription `json:"-"`

	OrderID uuid.UUID `json:"order_id" binding:"required" gorm:"type:uuid;not null;index"`
	Order   Order     `json:"-"`

	Amount   int    `json:"amount" binding:"required,gt=0" gorm:"not null;index"`
	Currency string `json:"currency" binding:"required,len=3" gorm:"size:3;not null;index"`

	PaidAt time.Time `json:"paid_at" gorm:"index"`

	PaymentStatus PaymentStatus `json:"payment_status" binding:"required" gorm:"size:20;not null;index"`
	Provider      string        `json:"provider" binding:"required,min=2,max=50" gorm:"size:50;not null;index"`
}
