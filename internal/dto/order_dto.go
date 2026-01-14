package dto

import "github.com/google/uuid"

type OrderCreateRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	ServiceID uuid.UUID `json:"service_id" binding:"required"`
	IsPaid    bool      `json:"is_paid"`
}

type OrderUpdateRequest struct {
	IsPaid *bool `json:"is_paid" binding:"omitempty"`
}
