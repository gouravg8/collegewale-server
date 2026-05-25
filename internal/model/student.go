package model

import (
	"collegeWaleServer/internal/enums"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	// Basic info
	FirstName string `gorm:"size:80;not null"`
	LastName  string `gorm:"size:80;not null"`
	Email     string `gorm:"uniqueIndex;size:100;not null"`
	Phone     string `gorm:"size:20"`
	// Academic info
	RollNumber       string `gorm:"not null"`
	CourseID         uint
	Course           Courses `gorm:"foreignKey:CourseID;references:ID;"`
	Year             int     `gorm:"not null"`
	Gender           string
	Semester         string
	EnrollmentNumber string
	//relation
	UserID        uint
	User          *User               `gorm:"foreignKey:UserID;references:ID;"`
	StudentStatus enums.StudentStatus `gorm:"default:'draft';not null"`
}
