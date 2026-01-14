package dto

type Pages struct {
	Page     int `json:"page" binding:"omitempty,min=1"`
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100"`
}
