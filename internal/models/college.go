package models

import (
	"collegeWaleServer/internal/enums"
	"time"

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
	Status     enums.CollegeType `gorm:"not null; defualt:'pending'"`

	// invicte based login
	InviteToken  string `gorm:"size:255;index"`
	InviteExpiry time.Time

	Subject []Subject
}
