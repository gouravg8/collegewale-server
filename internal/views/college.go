package views

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/college"
	"collegeWaleServer/internal/utils"
	"strings"
)

type College struct {
	Name       string             `json:"name"`
	Code       string             `json:"code"`
	Phone      string             `json:"phone"`
	Email      string             `json:"email"`
	CourseType college.CourseType `json:"course_type"` //TODO one college can have multiple courses
	Seats      uint               `json:"seats"`
	Logo       string             `json:"logo"`
}

func (c *College) IsValidRequest() error {
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
	if string(c.CourseType) == "" {
		return errz.NewBadRequest("course type cannot be empty")
	}
	if err := c.CourseType.IsValidCourseType(); err != nil {
		return err
	}
	if c.Seats <= 0 {
		return errz.NewBadRequest("seats must be greater than zero")
	}
	return nil
}
