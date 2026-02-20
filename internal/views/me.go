package views

import (
	"collegeWaleServer/internal/enums/roles"
)

type MeLogin struct {
	Username *string `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

type MeResponse struct {
	Email       string        `json:"email"`
	Phone       string        `json:"phone,omitempty"`
	Roles       []roles.Roles `json:"roles"`
	CollegeCode string        `json:"college_code,omitempty"`
	Token       string        `json:"token"`
}
