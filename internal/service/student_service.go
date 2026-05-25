package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums"
	"collegeWaleServer/internal/utils/common"
	views "collegeWaleServer/internal/views"
	v "collegeWaleServer/internal/views/common"
	"errors"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"

	"collegeWaleServer/internal/model"
)

var ErrStudentNotFound = errors.New("student not found")

// StudentService is the implementation for student-related operations
type StudentService struct {
	db *gorm.DB
}

// NewStudentService creates a new StudentService instance
func NewStudentService(DB *gorm.DB) *StudentService {
	return &StudentService{
		db: DB,
	}
}

// GetByID retrieves a student by ID
func (s *StudentService) GetByID(id int64) (*model.Student, error) {
	var student model.Student
	err := s.db.First(&student, uint(id))
	if err != nil {
		return nil, ErrStudentNotFound
	}

	return &student, nil
}

func (s *StudentService) ListStudents(filter views.StudentFilter) (v.DataList, error) {
	limit := 10
	if filter.PageSize > 0 {
		limit = filter.PageSize
	}

	offset := filter.PageNum
	courseID := filter.CourseID
	userID := filter.UserID
	rollNumber := filter.RollNumber
	gender := filter.Gender
	semester := filter.Semester
	year := filter.Year
	enrollment := filter.Enrollment

	query := s.db.Model(&model.Student{})

	if courseID != nil {
		query = query.Preload("Course")
	}

	if userID != nil {
		query = query.Preload("User")
	}

	if rollNumber != "" {
		query = query.Where("roll_number = ?", rollNumber)
	}

	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	if semester != "" {
		query = query.Where("semester = ?", semester)
	}

	if year != 0 {
		query = query.Where("year = ?", year)
	}

	if enrollment != "" {
		query = query.Where("enrollment_number = ?", enrollment)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return v.DataList{}, err
	}

	var students []model.Student
	err = query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		return v.DataList{}, err
	}
	if len(students) == 0 {
		return v.DataList{}, nil
	}
	myStudents := make([]any, 0)
	for _, student := range students {
		myStudents = append(myStudents, views.NewStudentInfoResponse(student))
	}
	response := v.NewAllDataList(myStudents)
	return response, nil
}

func (s *StudentService) UpdateStudentStatus(mid common.MaskedId, status enums.StudentStatus) error {
	id := common.Unmask(mid)
	if result := s.db.Model(&model.Student{}).
		Where("id = ? AND status IN (?)", id, status.GetAllowedStatus()).
		Update("status", status); result.Error != nil {
		log.Errorf("update student status error: %v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errz.NewBadRequest("student not found")
		}
		return errz.NewBadRequest("failed to update student status.")
	}
	return nil
}
