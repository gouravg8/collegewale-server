package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	StudentID uint
	Student   Student
	SubjectID uint
	Subject   Subject
	Date      time.Time
	Status    string `gorm:"size:20"` // Present / Absent / Late
}
