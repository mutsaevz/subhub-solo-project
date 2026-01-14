package dto

// DTO для категории сервиса
type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type CategoryUpdateRequest struct {
	Name *string `json:"name" binding:"omitempty,min=2,max=100"`
}
