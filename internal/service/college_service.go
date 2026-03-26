package service

import (
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/views"

	"gorm.io/gorm"
)

type CollegeService struct {
	db *gorm.DB
}

func NewCollegeService(db *gorm.DB) *CollegeService {
	return &CollegeService{db: db}
}

func (s *CollegeService) GetStats(user *model.User) (views.CollegeStatsResponse, error) {
	var stats struct {
		Total    int64
		Approved int64
		Pending  int64
		Rejected int64
		Unpayed  int64
	}
	if err := s.db.Model(&model.Student{}).Joins(`JOIN "user" ON "user".id = student.user_id`).Select(`
            COUNT(CASE WHEN student.student_status = 'approved' THEN 1 END) as approved,
            COUNT(CASE WHEN student.student_status = 'rejected' THEN 1 END) as rejected,
            COUNT(CASE WHEN student.student_status = 'pending' THEN 1 END) as pending,
            COUNT(CASE WHEN student.student_status = 'unpayed' THEN 1 END) as unpayed,
            COUNT(*) as total
        `).Where(`"user".college_id = ?`, user.CollegeID).
		Scan(&stats).Error; err != nil {
		return views.CollegeStatsResponse{}, err
	}
	// These counts would typically come from an Applications model which is not fully shown,
	// but following the structure of the requested view:
	return views.CollegeStatsResponse{
		CollegeName:          user.College.Name,
		TotalStudents:        stats.Total,
		TotalApplications:    stats.Pending,
		TotalApproved:        stats.Approved,
		TotalPendingPayments: stats.Unpayed,
	}, nil

}
