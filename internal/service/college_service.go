package service

import (
	"collegeWaleServer/internal/enums"
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
		Total     int64
		Draft     int64
		Submitted int64
		Verified  int64
		Admitted  int64
		Approved  int64
		Unpayed   int64
	}
	if err := s.db.Model(&model.Student{}).Joins(`JOIN "user" ON "user".id = student.user_id`).Select(`
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as draft,
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as submitted,
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as verified,
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as admitted,
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as approved,
            COUNT(CASE WHEN student.student_status = ? THEN 1 END) as unpayed,
            COUNT(*) as total
        `, enums.Draft, enums.Submitted, enums.Verified, enums.Admitted, enums.Approved, enums.Unpayed).Where(`"user".college_id = ?`, user.CollegeID).
		Scan(&stats).Error; err != nil {
		return views.CollegeStatsResponse{}, err
	}
	// These counts would typically come from an Applications model which is not fully shown,
	// but following the structure of the requested view:
	return views.CollegeStatsResponse{
		CollegeName:          user.College.Name,
		TotalStudents:        stats.Total,
		TotalDraft:           stats.Draft,
		TotalSubmitted:       stats.Submitted,
		TotalVerified:        stats.Verified,
		TotalApproved:        stats.Approved,
		TotalPendingPayments: stats.Unpayed,
	}, nil

}
