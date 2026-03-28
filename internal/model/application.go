package model

import (
	"collegeWaleServer/internal/enums/application"
	"collegeWaleServer/internal/enums/college"
	"time"

	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	StudentID    uint                          `gorm:"not null;index"`
	CollegeID    uint                          `gorm:"not null;index"`
	CourseType   college.CourseType            `gorm:"not null"`
	AcademicYear string                        `gorm:"size:20;not null"`
	Status       application.ApplicationStatus `gorm:"size:20;not null;default:'draft';index"`
	Remarks      string                        `gorm:"type:text"`

	// Transition timestamps
	SubmittedAt *time.Time
	VerifiedAt  *time.Time
	ApprovedAt  *time.Time
	AdmittedAt  *time.Time
	RejectedAt  *time.Time

	// Who performed each transition
	SubmittedBy *uint
	VerifiedBy  *uint
	ApprovedBy  *uint
	AdmittedBy  *uint

	// Relations
	Student Student  `gorm:"foreignKey:StudentID;references:ID"`
	College College  `gorm:"foreignKey:CollegeID;references:ID"`
}
