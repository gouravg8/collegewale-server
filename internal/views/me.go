package views

import "collegeWaleServer/internal/enums"

type MeLogin struct {
	Username *string `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

type Me struct {
	Email     string        `json:"email"`
	Phone     string        `json:"phone;omitempty"`
	Role      []enums.Roles `json:"role"`
	CollegeID uint          `json:"college_id;omitempty"`
}
