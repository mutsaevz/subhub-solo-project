package models

// Категория сервиса
type Category struct {
	Base

	Name string `json:"name" binding:"required,min=2,max=100" gorm:"size:100;not null"`

	Services []Service `json:"-" gorm:"foreignKey:CategoryID"`
}
