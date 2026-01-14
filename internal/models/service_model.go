package models

import "github.com/google/uuid"

// Сервис, на который можно оформить подписку
type Service struct {
	Base

	Name string `json:"name" binding:"required,min=2,max=100" gorm:"size:100;not null;index"`

	CategoryID uuid.UUID `json:"category_id" binding:"required" gorm:"type:uuid;not null;index"`
	Category   Category  `json:"-"`

	Website string `json:"website" binding:"omitempty,url" gorm:"size:255;index"`
	LogoUrl string `json:"logo_url" binding:"omitempty,url" gorm:"size:255"`

	Subscriptions []Subscription `json:"-" gorm:"foreignKey:ServiceID"`
}
