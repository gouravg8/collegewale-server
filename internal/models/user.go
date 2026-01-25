package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"` // "student", "college_admin", "admin"
	CollegeID    *uint
}
