package views

import (
	"collegeWaleServer/internal/enums/college"
)

type CollegeSignup struct {
	Name       string             `json:"name"`
	Code       string             `json:"code"`
	Phone      string             `json:"phone"`
	Email      string             `json:"email"`
	CourseType college.CourseType `json:"course_type"`
	Seats      uint               `json:"seats"`
	Logo       string             `json:"logo"`
}
