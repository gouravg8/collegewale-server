package views

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/utils"
	"strings"
)

type CollegeRequest struct {
	Name    string   `json:"name"`
	Code    string   `json:"code"`
	Phone   string   `json:"phone"`
	Email   string   `json:"email"`
	Courses []string `json:"available_courses"`
	Seats   uint     `json:"seats"`
}

type CollegeResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Logo string `json:"logo"`
}

func (c *CollegeRequest) IsValidRequest() error {
	if strings.TrimSpace(c.Name) == "" {
		return errz.NewBadRequest("college name cannot be empty")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errz.NewBadRequest("email cannot be empty")
	}
	if !utils.IsEmailValid(c.Email) {
		return errz.NewBadRequest("invalid email format")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errz.NewBadRequest("phone cannot be empty")
	}
	if !utils.IsPhoneValid(c.Phone) {
		return errz.NewBadRequest("invalid phone format")
	}
	if strings.TrimSpace(c.Code) == "" {
		return errz.NewBadRequest("college code cannot be empty")
	}
	if len(c.Courses) == 0 {
		return errz.NewBadRequest("college should provide at least one course")
	}
	if c.Seats <= 0 {
		return errz.NewBadRequest("seats must be greater than zero")
	}
	return nil
}

type CollegeSignup struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

func (c CollegeSignup) IsValid() error {
	if strings.TrimSpace(c.Password) == "" {
		return errz.NewBadRequest("password is required")
	}
	if strings.TrimSpace(c.Username) == "" {
		return errz.NewBadRequest("username is required")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errz.NewBadRequest("email is required")
	}
	if !utils.IsPhoneValid(strings.TrimSpace(c.Phone)) {
		return errz.NewBadRequest("phone is required")
	}
	return nil
}

// CollegeWithAdminRequest is the minimal payload for the combined
// "create college + its college admin" flow: just the college name and
// the admin's login details. Everything else (code, seats, phone/email
// fallbacks) is derived or defaulted server-side.
type CollegeWithAdminRequest struct {
	CollegeName   string `json:"college_name"`
	AdminUsername string `json:"admin_username"`
	AdminEmail    string `json:"admin_email"`
	AdminPhone    string `json:"admin_phone"`
	AdminPassword string `json:"admin_password"`
}

func (c CollegeWithAdminRequest) IsValid() error {
	if strings.TrimSpace(c.CollegeName) == "" {
		return errz.NewBadRequest("college name is required")
	}
	if strings.TrimSpace(c.AdminUsername) == "" {
		return errz.NewBadRequest("admin username is required")
	}
	if strings.TrimSpace(c.AdminEmail) == "" {
		return errz.NewBadRequest("admin email is required")
	}
	if !utils.IsEmailValid(c.AdminEmail) {
		return errz.NewBadRequest("invalid admin email format")
	}
	if strings.TrimSpace(c.AdminPassword) == "" {
		return errz.NewBadRequest("admin password is required")
	}
	return nil
}

type CollegeStatsResponse struct {
	CollegeName          string `json:"college_name"`
	TotalStudents        int64  `json:"total_students"`
	TotalDraft           int64  `json:"total_draft"`
	TotalSubmitted       int64  `json:"total_submitted"`
	TotalVerified        int64  `json:"total_verified"`
	TotalApproved        int64  `json:"total_approved"`
	TotalAdmitted        int64  `json:"total_admitted"`
	TotalPendingPayments int64  `json:"total_pending_payments"`
}
