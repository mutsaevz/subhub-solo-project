package models

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Пользователь системы
type User struct {
	Base

	Email string `json:"email" binding:"required,email" gorm:"size:255;not null;uniqueIndex"`

	Password string `json:"-" gorm:"not null"`

	FirstName string `json:"first_name" binding:"required,min=2,max=100" gorm:"size:100;not null;index"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100" gorm:"size:100;not null;index"`

	Roles         Role           `json:"role" binding:"required,oneof=user admin" gorm:"type:text;not null;default:'user';index;check:roles IN ('user','admin')"`
	Subscriptions []Subscription `json:"-" gorm:"foreignKey:UserID"`
}
