package views

import "collegeWaleServer/internal/enums/college"

type StudentInfoResponse struct { //TODO add more info as required per user
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	RollNumber string             `json:"roll_number"`
	CourseType college.CourseType `json:"course_type"`
	Year       int                `json:"year"`
	Subjects   []string           `json:"subjects"`
	//// relationships
	//CollegeID uint    `gorm:"not null"` // fk -> college
	//College   College `gorm:"foreignKey:CollegeID"`

	//Subject []Subject `gorm:"many2many:student_subjects"`
}
