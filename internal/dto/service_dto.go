package dto

import "github.com/google/uuid"

// DTO для создания и обновления сервиса
type ServiceCreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`

	CategoryID uuid.UUID `json:"category_id" binding:"required"`

	Website string `json:"website" binding:"omitempty,url"`
	LogoUrl string `json:"logo_url" binding:"omitempty,url"`
}

type ServiceUpdateRequest struct {
	Name *string `json:"name" binding:"omitempty,min=2,max=100"`

	CategoryID *uuid.UUID `json:"category_id" binding:"omitempty"`

	Website *string `json:"website" binding:"omitempty,url"`
	LogoUrl *string `json:"logo_url" binding:"omitempty,url"`
}
