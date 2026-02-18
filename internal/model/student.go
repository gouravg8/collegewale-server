package model

import (
	"collegeWaleServer/internal/enums/college"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	// Basic info
	FirstName string `gorm:"size:80;not null"`
	LastName  string `gorm:"size:80;not null"`
	Email     string `gorm:"uniqueIndex;varchar(255);not null"`
	Phone     string `gorm:"size:20"`
	// Academic info
	RollNumber       string             `gorm:"not null"`
	CourseType       college.CourseType `gorm:"not null"`
	Year             int                `gorm:"not null"`
	Gender           string
	Semester         string
	EnrollmentNumber string
	CollegeCode      string `gorm:"notnull"` //can be used to get college info

	Subject []Subject `gorm:"many2many:student_subjects"`
}
