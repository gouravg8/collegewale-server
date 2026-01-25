package models

import (
	"collegeWaleServer/internal/enums"

	"gorm.io/gorm"
)

type College struct {
	gorm.Model
	Name       string           `gorm:"type:text"`
	Code       string           `gorm:"not null"`
	Phone      string           `gorm:"not null"`
	Email      string           `gorm:"not null"`
	CourseType enums.CourseType `gorm:"not null; default:'GNM'"`
	Seats      uint             `gorm:"not null"`
	Logo       string
}
