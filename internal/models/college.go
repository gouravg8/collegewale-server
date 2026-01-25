package models

import (
	"collegeWaleServer/internal/enums"

	"gorm.io/gorm"
)

type College struct {
	gorm.Model
	Name       string           `gorm:"type:text;unique"`
	Code       string           `gorm:"not null;unique"`
	Phone      string           `gorm:"not null;unique"`
	Email      string           `gorm:"not null;unique"`
	CourseType enums.CourseType `gorm:"not null; default:'GNM'"`
	Seats      uint             `gorm:"not null"`
	Logo       string
}
