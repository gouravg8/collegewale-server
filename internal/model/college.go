package model

import (
	"time"

	"gorm.io/gorm"
)

type College struct {
	gorm.Model
	Name         string `gorm:"type:text;unique"`
	Code         string `gorm:"not null;unique"`
	Phone        string `gorm:"not null;unique"`
	Email        string `gorm:"not null;unique"`
	Seats        uint   `gorm:"not null"`
	Logo         string
	Courses      []Courses `gorm:"many2many:college_courses;"`
	PasswordHash string
	CreatedById  uint

	// invite based login
	InviteToken  string `gorm:"type:text"`
	InviteExpiry time.Time
}
