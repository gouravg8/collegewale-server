package application

import (
	"collegeWaleServer/errz"
	appEnum "collegeWaleServer/internal/enums/application"
	"collegeWaleServer/internal/enums/college"
	"collegeWaleServer/internal/model"
	"strings"
	"time"
)

// --- Requests ---

type CreateApplicationRequest struct {
	StudentID    uint               `json:"student_id"`
	CourseType   college.CourseType `json:"course_type"`
	AcademicYear string             `json:"academic_year"`
	Remarks      string             `json:"remarks"`
}

func (r CreateApplicationRequest) IsValid() error {
	if r.StudentID == 0 {
		return errz.NewBadRequest("student_id is required")
	}
	if err := r.CourseType.IsValidCourseType(); err != nil {
		return err
	}
	if strings.TrimSpace(r.AcademicYear) == "" {
		return errz.NewBadRequest("academic_year is required")
	}
	return nil
}

type UpdateStatusRequest struct {
	Status  appEnum.ApplicationStatus `json:"status"`
	Remarks string                    `json:"remarks"`
}

func (r UpdateStatusRequest) IsValid() error {
	return r.Status.IsValid()
}

// --- Responses ---

type ApplicationResponse struct {
	ID           uint                      `json:"id"`
	StudentID    uint                      `json:"student_id"`
	StudentName  string                    `json:"student_name"`
	CollegeID    uint                      `json:"college_id"`
	CollegeName  string                    `json:"college_name"`
	CourseType   college.CourseType        `json:"course_type"`
	AcademicYear string                    `json:"academic_year"`
	Status       appEnum.ApplicationStatus `json:"status"`
	Remarks      string                    `json:"remarks"`
	SubmittedAt  *time.Time                `json:"submitted_at,omitempty"`
	VerifiedAt   *time.Time                `json:"verified_at,omitempty"`
	ApprovedAt   *time.Time                `json:"approved_at,omitempty"`
	AdmittedAt   *time.Time                `json:"admitted_at,omitempty"`
	RejectedAt   *time.Time                `json:"rejected_at,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

func NewApplicationResponse(a *model.Application) ApplicationResponse {
	resp := ApplicationResponse{
		ID:           a.ID,
		StudentID:    a.StudentID,
		CollegeID:    a.CollegeID,
		CourseType:   a.CourseType,
		AcademicYear: a.AcademicYear,
		Status:       a.Status,
		Remarks:      a.Remarks,
		SubmittedAt:  a.SubmittedAt,
		VerifiedAt:   a.VerifiedAt,
		ApprovedAt:   a.ApprovedAt,
		AdmittedAt:   a.AdmittedAt,
		RejectedAt:   a.RejectedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.Student.ID != 0 {
		resp.StudentName = a.Student.FirstName + " " + a.Student.LastName
	}
	if a.College.ID != 0 {
		resp.CollegeName = a.College.Name
	}
	return resp
}
