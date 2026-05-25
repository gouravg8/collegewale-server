package model

import "gorm.io/gorm"

type Courses struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"type:text;default:''"`
	Subjects    []Subject `gorm:"many2many:course_subjects;"`
}

type CollegeCourses struct {
	CollegeId uint `gorm:"not null"`
	CourseId  uint `gorm:"not null"`
}

type CourseSubjects struct {
	gorm.Model
	CourseId  uint `gorm:"not null"`
	SubjectId uint `gorm:"not null"`
}
