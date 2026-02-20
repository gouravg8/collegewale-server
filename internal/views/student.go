package views

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums/college"
	"collegeWaleServer/internal/utils"
	"strings"
)

type StudentInfoResponse struct {
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	RollNumber       string             `json:"roll_number"`
	CourseType       college.CourseType `json:"course_type"`
	Year             int                `json:"year"`
	Subjects         []string           `json:"subjects"`
	EnrollmentNumber string             `json:"enrollment_number"`
	Semester         string             `json:"semester"`
	Gender           string             `json:"gender"`
}

type StudentForm struct { //TODO add more info as required per user
	Username         string             `json:"username"`
	Password         string             `json:"password"`
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	Email            string             `json:"email"`
	Phone            string             `json:"phone"`
	RollNumber       string             `json:"roll_number"`
	CourseType       college.CourseType `json:"course_type"`
	Year             int                `json:"year"`
	Gender           string             `json:"gender"`
	Semester         string             `json:"semester"`
	EnrollmentNumber string             `json:"enrollment_number"`
	Subjects         []string           `json:"subjects"`
}

func (s StudentForm) IsValid() error {
	if strings.TrimSpace(s.Password) == "" {
		return errz.NewBadRequest("password is required")
	}
	if strings.TrimSpace(s.Username) == "" {
		return errz.NewBadRequest("username is required")
	}
	if strings.TrimSpace(s.FirstName) == "" {
		return errz.NewBadRequest("first name is required")
	}
	if strings.TrimSpace(s.LastName) == "" {
		return errz.NewBadRequest("last name is required")
	}
	if !utils.IsEmailValid(s.Email) {
		return errz.NewBadRequest("email is required")
	}
	if !utils.IsPhoneValid(s.Phone) {
		return errz.NewBadRequest("phone is required")
	}
	if strings.TrimSpace(s.RollNumber) == "" {
		return errz.NewBadRequest("role number is required")
	}
	err := s.CourseType.IsValidCourseType()
	if err != nil {
		return err
	}
	if s.Year <= 0 {
		return errz.NewBadRequest("year cannot be zero")
	}
	if strings.TrimSpace(s.Gender) == "" {
		return errz.NewBadRequest("gender is required")
	}
	if strings.TrimSpace(s.Semester) == "" {
		return errz.NewBadRequest("semester is required")
	}
	if len(s.Subjects) == 0 {
		return errz.NewBadRequest("please provide some subjects")
	}

	return nil
}
