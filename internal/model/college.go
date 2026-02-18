package model

import (
	"collegeWaleServer/internal/enums"
	"collegeWaleServer/internal/enums/college"
	"time"

	"gorm.io/gorm"
)

type College struct {
	gorm.Model
	Name         string             `gorm:"type:text;unique"`
	Code         string             `gorm:"not null;unique"`
	Phone        string             `gorm:"not null;unique"`
	Email        string             `gorm:"not null;unique"`
	CourseType   college.CourseType `gorm:"not null; default:'gnm'"`
	Seats        uint               `gorm:"not null"`
	Logo         string
	Status       enums.CollegeType `gorm:"not null; defualt:'pending'"`
	PasswordHash string
	CreatedById  uint

	// invite based login
	InviteToken  string `gorm:"type:text"`
	InviteExpiry time.Time

	Subjects []Subject
}
