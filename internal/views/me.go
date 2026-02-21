package views

import (
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
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

type MyInfo struct {
	Email    string               `json:"email"`
	Username string               `json:"username"`
	Phone    string               `json:"phone"`
	Roles    []roles.Roles        `json:"roles"`
	Student  *StudentInfoResponse `json:"student_info,omitempty"`
	College  *CollegeResponse     `json:"College_info,omitempty"`
}

func NewMyInfo(u model.User) MyInfo {
	myRoles := make([]roles.Roles, 0)
	for _, r := range u.Roles {
		myRoles = append(myRoles, r.Name)
	}
	myInfo := MyInfo{
		Email:    u.Email,
		Username: u.Username,
		Roles:    myRoles,
	}
	if u.Phone != nil {
		myInfo.Phone = *u.Phone
	}
	if u.Student != nil {
		s := u.Student
		subjects := make([]string, 0)
		for _, sbj := range s.Subject {
			subjects = append(subjects, sbj.Name)
		}
		myInfo.Student = &StudentInfoResponse{
			FirstName:        s.FirstName,
			LastName:         s.LastName,
			RollNumber:       s.RollNumber,
			CourseType:       s.CourseType,
			Year:             s.Year,
			Subjects:         subjects,
			EnrollmentNumber: s.EnrollmentNumber,
			Semester:         s.Semester,
			Gender:           s.Gender,
		}
	}
	if u.College != nil {
		c := u.College
		myInfo.College = &CollegeResponse{
			Name: c.Name,
			Code: c.Code,
			Logo: c.Logo,
		}
	}

	return myInfo
}
