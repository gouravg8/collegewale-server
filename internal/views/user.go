package views

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/utils"
	"strings"
)

type UpdateUserRequest struct {
	Email       *string        `json:"email,omitempty"`
	Username    *string        `json:"username,omitempty"`
	Phone       *string        `json:"phone,omitempty"`
	Password    *string        `json:"password,omitempty"`
	CollegeCode *string        `json:"college_code,omitempty"`
	Roles       *[]roles.Roles `json:"roles,omitempty"`
}

func (r UpdateUserRequest) IsValid() error {
	if r.Email == nil && r.Username == nil && r.Phone == nil && r.Password == nil && r.CollegeCode == nil && r.Roles == nil {
		return errz.NewBadRequest("no fields to update")
	}

	if r.Email != nil {
		e := strings.TrimSpace(*r.Email)
		if e == "" || !utils.IsEmailValid(e) {
			return errz.NewBadRequest("invalid email format")
		}
	}

	if r.Username != nil {
		if strings.TrimSpace(*r.Username) == "" {
			return errz.NewBadRequest("username cannot be empty")
		}
	}

	if r.Phone != nil {
		p := strings.TrimSpace(*r.Phone)
		if p == "" || !utils.IsPhoneValid(p) {
			return errz.NewBadRequest("invalid phone format")
		}
	}

	if r.Password != nil {
		p := strings.TrimSpace(*r.Password)
		if p == "" {
			return errz.NewBadRequest("password cannot be empty")
		}
		if len(p) < 6 {
			return errz.NewBadRequest("password must be at least 6 characters")
		}
	}

	if r.Roles != nil {
		if len(*r.Roles) == 0 {
			return errz.NewBadRequest("roles cannot be empty")
		}
		// Validate role values against allowed constants
		for _, rr := range *r.Roles {
			switch rr {
			case roles.Admin, roles.Student, roles.College:
				// valid
			default:
				return errz.NewBadRequest("invalid role provided")
			}
		}
	}

	// CollegeCode if present should not be blank
	if r.CollegeCode != nil && strings.TrimSpace(*r.CollegeCode) == "" {
		return errz.NewBadRequest("college code cannot be empty")
	}

	return nil
}
