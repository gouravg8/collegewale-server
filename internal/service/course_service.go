package service

import (
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/views"
	view "collegeWaleServer/internal/views/common"

	"gorm.io/gorm"
)

type CourseService struct {
	db *gorm.DB
}

func NewCourseService(DB *gorm.DB) *CourseService {
	return &CourseService{DB}
}

func (cs *CourseService) ListCourses(filter views.CoursesFilter) (view.DataList, error) {
	// Fetch courses with name filter (case-insensitive)
	var courses []model.Courses
	query := cs.db.Model(&model.Courses{})
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.WithSubjects {
		query = query.Preload("Subjects")
	}
	if err := query.Find(&courses).Error; err != nil {
		return view.DataList{}, err
	}
	data := make([]any, 0)
	for _, c := range courses {
		data = append(data, views.NewCoursesResponse(c))
	}
	return view.NewAllDataList(data), nil
}
