package service

import (
	views "collegeWaleServer/internal/views"
	v "collegeWaleServer/internal/views/common"
	"context"
	"errors"

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
func (s *StudentService) GetByID(ctx context.Context, id int64) (*model.Student, error) {
	var student model.Student
	err := s.db.First(&student, uint(id))
	if err != nil {
		return nil, ErrStudentNotFound
	}

	return &student, nil
}

// Create creates a new student record
func (s *StudentService) Create(ctx context.Context, student *model.Student) (*model.Student, error) {
	var created model.Student
	err := s.db.Create(&student).Error
	if err != nil {
		return nil, err
	}

	created = *student
	return &created, nil
}

// Update updates an existing student record
func (s *StudentService) Update(ctx context.Context, id int64, updates map[string]any) (*model.Student, error) {
	var student model.Student
	err := s.db.First(&student, uint(id)).Error
	if err != nil {
		return nil, ErrStudentNotFound
	}

	for key, value := range updates {
		switch key {
		case "course":
			if courseID, ok := value.(uint); ok {
				var course model.Courses
				err := s.db.First(&course, courseID).Error
				if err != nil {
					return nil, err
				}
				student.CourseID = courseID
			}
		case "year":
			if yearVal, ok := value.(int); ok {
				student.Year = yearVal
			}
		case "gender":
			if genderVal, ok := value.(string); ok {
				student.Gender = genderVal
			}
		case "semester":
			if semesterVal, ok := value.(string); ok {
				student.Semester = semesterVal
			}
		case "enrollment_number":
			if enrollmentVal, ok := value.(string); ok {
				student.EnrollmentNumber = enrollmentVal
			}
		}
	}

	err = s.db.Save(&student).Error
	if err != nil {
		return nil, err
	}

	var updated model.Student
	s.db.First(&updated, uint(id))
	return &updated, nil
}

// Delete deletes a student record
func (s *StudentService) Delete(ctx context.Context, id int64) error {
	err := s.db.Delete(&model.Student{}, uint(id))
	if err != nil {
		return ErrStudentNotFound
	}

	return nil
}
func (s *StudentService) ListStudents(user *model.User, filter views.StudentFilter) (v.DataList, error) {
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
