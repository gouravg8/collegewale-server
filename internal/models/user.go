package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Username     string `gorm:"uniqueIndex;not null"`
	Phone        *string
	PasswordHash string `gorm:"type:text;not null"`
	Role         []Role `gorm:"many2many:user_roles;"`
	CollegeID    *uint
}
