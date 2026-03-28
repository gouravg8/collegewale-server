package application

import (
	"collegeWaleServer/errz"
	appEnum "collegeWaleServer/internal/enums/application"
	"collegeWaleServer/internal/model"
	appViews "collegeWaleServer/internal/views/application"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type ApplicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

func (s *ApplicationService) CreateApplication(req appViews.CreateApplicationRequest, user *model.User) (*model.Application, error) {
	if user.CollegeID == nil {
		return nil, errz.NewBadRequest("user is not associated with a college")
	}

	// Verify student exists and belongs to the same college
	var student model.Student
	if err := s.db.First(&student, req.StudentID).Error; err != nil {
		return nil, errz.NewNotFound("student not found")
	}
	if student.CollegeCode != user.College.Code {
		return nil, errz.NewForbidden("student does not belong to your college")
	}

	app := model.Application{
		StudentID:    req.StudentID,
		CollegeID:    *user.CollegeID,
		CourseType:   req.CourseType,
		AcademicYear: strings.TrimSpace(req.AcademicYear),
		Status:       appEnum.Draft,
		Remarks:      strings.TrimSpace(req.Remarks),
	}

	if err := s.db.Create(&app).Error; err != nil {
		log.Errorf("failed to create application: %v", err)
		return nil, errz.NewBadRequest("failed to create application")
	}

	// Reload with associations
	if err := s.db.Preload("Student").Preload("College").First(&app, app.ID).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *ApplicationService) GetApplication(id uint) (*model.Application, error) {
	var app model.Application
	if err := s.db.Preload("Student").Preload("College").First(&app, id).Error; err != nil {
		return nil, errz.NewNotFound(fmt.Sprintf("application %d not found", id))
	}
	return &app, nil
}

func (s *ApplicationService) ListApplications(collegeID uint, status string, page, pageSize int) ([]model.Application, int64, error) {
	query := s.db.Model(&model.Application{}).Where("college_id = ?", collegeID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var apps []model.Application
	err := query.
		Preload("Student").Preload("College").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}
	return apps, total, nil
}

func (s *ApplicationService) UpdateStatus(id uint, req appViews.UpdateStatusRequest, user *model.User) (*model.Application, error) {
	var app model.Application
	if err := s.db.Preload("Student").Preload("College").First(&app, id).Error; err != nil {
		return nil, errz.NewNotFound(fmt.Sprintf("application %d not found", id))
	}

	// Enforce state machine
	if err := app.Status.CanTransitionTo(req.Status); err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]any{
		"status":  req.Status,
		"remarks": strings.TrimSpace(req.Remarks),
	}

	switch req.Status {
	case appEnum.Submitted:
		updates["submitted_at"] = &now
		updates["submitted_by"] = &user.ID
	case appEnum.Verified:
		updates["verified_at"] = &now
		updates["verified_by"] = &user.ID
	case appEnum.Approved:
		updates["approved_at"] = &now
		updates["approved_by"] = &user.ID
	case appEnum.Admitted:
		updates["admitted_at"] = &now
		updates["admitted_by"] = &user.ID
	case appEnum.Rejected:
		updates["rejected_at"] = &now
	case appEnum.Draft:
		// Reset timestamps when going back to draft from rejected
		updates["submitted_at"] = nil
		updates["verified_at"] = nil
		updates["approved_at"] = nil
		updates["admitted_at"] = nil
		updates["rejected_at"] = nil
		updates["submitted_by"] = nil
		updates["verified_by"] = nil
		updates["approved_by"] = nil
		updates["admitted_by"] = nil
	}

	// Use transaction for admitted — also decrement seat count
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Application{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if req.Status == appEnum.Admitted {
			result := tx.Model(&model.College{}).
				Where("id = ? AND seats > 0", app.CollegeID).
				Update("seats", gorm.Expr("seats - 1"))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errz.NewBadRequest("no available seats in this college")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Reload
	if err := s.db.Preload("Student").Preload("College").First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *ApplicationService) GetStudentApplications(studentID uint) ([]model.Application, error) {
	var apps []model.Application
	err := s.db.
		Preload("College").Preload("Student").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}
