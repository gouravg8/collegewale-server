package auth_view

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/utils"
	"strings"
)

type CollegeSignup struct {
	Name    string   `json:"name"`
	Code    string   `json:"code"`
	Phone   string   `json:"phone"`
	Email   string   `json:"email"`
	Courses []string `json:"courses"`
	Seats   uint     `json:"seats"`
	Logo    string   `json:"logo"`
}

func (cs CollegeSignup) IsValid() error {
	// --- Input Validation ---
	if strings.TrimSpace(cs.Name) == "" {
		return errz.NewBadRequest("college name cannot be empty")
	}
	if strings.TrimSpace(cs.Email) == "" {
		return errz.NewBadRequest("email cannot be empty")
	}
	if !utils.IsEmailValid(cs.Email) {
		return errz.NewBadRequest("invalid email format")
	}
	if strings.TrimSpace(cs.Phone) == "" {
		return errz.NewBadRequest("phone cannot be empty")
	}
	if !utils.IsPhoneValid(cs.Phone) {
		return errz.NewBadRequest("invalid phone format")
	}
	if strings.TrimSpace(cs.Code) == "" {
		return errz.NewBadRequest("college code cannot be empty")
	}
	if len(cs.Courses) == 0 {
		return errz.NewBadRequest("college should provide at least one course")
	}
	if cs.Seats <= 0 {
		return errz.NewBadRequest("seats must be greater than zero")
	}
	return nil
}

type SetPassword struct {
	Code            string `json:"code"`
	Email           string `json:"email"`
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
