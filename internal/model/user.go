package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Username     string `gorm:"uniqueIndex;not null"`
	Phone        *string
	PasswordHash string `gorm:"type:text;not null"`
	Roles        []Role `gorm:"many2many:user_roles;"`
	CollegeID    *uint
	College      *College `gorm:"foreignKey:CollegeID;references:ID;"`
	Student      *Student
	CreatedByID  uint
}

type UserRole struct {
	gorm.Model
	UserId uint `gorm:"not null"`
	RoleId uint `gorm:"not null"`
}
