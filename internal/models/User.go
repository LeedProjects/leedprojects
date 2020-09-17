package models

import (
	"time"

	"gorm.io/gorm"
)

//User model
type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Username  string          `gorm:"unique" json:"username"`
	Email     string          `gorm:"unique" json:"email"`
	Password  string          `json:"password,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
