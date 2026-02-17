package model

import (
	"collegeWaleServer/internal/enums/college"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	// Basic info
	FirstName string `gorm:"size:100;not null"`
	LastName  string `gorm:"size:100;not null"`
	Email     string `gorm:"uniqueIndex;size:255;not null"`
	Phone     string `gorm:"size:20"`

	UserID int
	User   User

	// Academic info
	RollNumber string             `gorm:"uniqueIndex;size:50;not null"`
	CourseType college.CourseType `gorm:"not null"`
	Year       int                `gorm:"not null"`

	// relationships
	CollegeID uint    `gorm:"not null"` // fk -> college
	College   College `gorm:"foreignKey:CollegeID"`

	Subject []Subject `gorm:"many2many:student_subjects"`
}
