package auth_view

import "collegeWaleServer/internal/enums"

type CollegeSignup struct {
	Name       string           `json:"name"`
	Code       string           `json:"code"`
	Phone      string           `json:"phone"`
	Email      string           `json:"email"`
	CourseType enums.CourseType `json:"course_type"`
	Seats      uint             `json:"seats"`
	Logo       string           `json:"logo"`
}
