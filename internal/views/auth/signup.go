package auth_view

import (
	"collegeWaleServer/internal/enums"
)

type CollegeSignup struct {
	Name       string           `json:"name"`
	Code       string           `json:"code"`
	Phone      string           `json:"phone"`
	Email      string           `json:"email"`
	CourseType enums.CourseType `json:"course_type"`
	Seats      uint             `json:"seats"`
	Logo       string           `json:"logo"`
}

type SetPassword struct {
	CollegeID       string `json:"college_id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type CollegeLogin struct {
	Code     string `json:"code"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CollegeLoginResponse struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Email string `json:"email"`
	Token string `json:"token"`
}
