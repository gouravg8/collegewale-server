package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/db"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/utils"
	"collegeWaleServer/internal/views"
	"errors"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type RegistryService struct {
	db *gorm.DB
}

func NewRegistryService(db *gorm.DB) *RegistryService {
	return &RegistryService{db}
}

func (s RegistryService) RegisterCollege(req views.College, user *model.User) error {
	clg := model.College{
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Phone:       strings.TrimSpace(req.Phone),
		Email:       strings.TrimSpace(req.Email),
		CourseType:  req.CourseType,
		Seats:       req.Seats,
		Logo:        req.Logo,
		CreatedById: user.ID,
	}

	if err := s.db.Model(&model.College{}).Create(&clg).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ex := pgErr.Detail
			switch {
			case strings.Contains(ex, "name"):
				return errz.NewBadRequest("college name already exists")
			case strings.Contains(ex, "email"):
				return errz.NewBadRequest("email  already exists")
			case strings.Contains(ex, "phone"):
				return errz.NewBadRequest("phone already linked with another college")
			case strings.Contains(ex, "code"):
				return errz.NewBadRequest("college code already exists")
			default:
				return errz.NewBadRequest("college already exists")
			}
		}
		return err
	}
	return nil
}

func (s RegistryService) RegisterStudent(req views.StudentForm, user *model.User) error {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Errorf("Failed to hash password: %v", err)
		return errz.NewBadRequest("failed to save user password")
	}
	var role model.Role
	err = s.db.Model(&model.Role{}).Where("name = ?", roles.Student).First(&role).Error
	if err != nil {
		log.Errorf("Failed to find student role: %v", err)
		return errz.NewBadRequest("role not found")
	}
	var student = model.Student{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Phone:            req.Phone,
		RollNumber:       req.RollNumber,
		CourseType:       req.CourseType,
		Year:             req.Year,
		Gender:           req.Gender,
		Semester:         req.Semester,
		EnrollmentNumber: req.EnrollmentNumber,
	}

	var me = model.User{
		Email:        strings.TrimSpace(req.Email),
		Username:     strings.TrimSpace(req.Username),
		PasswordHash: passwordHash,
		Roles:        []model.Role{role},
		CollegeID:    user.CollegeID,
		Student:      &student,
		CreatedByID:  user.ID,
	}
	if user.College != nil {
		student.CollegeCode = user.College.Code
	}
	cleanedPhone := strings.TrimSpace(req.Phone)
	if cleanedPhone != "" {
		me.Phone = &cleanedPhone
	}
	err = db.DB.Model(&model.User{}).Create(&me).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ex := pgErr.Detail
			switch {
			case strings.Contains(ex, "username"):
				return errz.NewBadRequest("username already exists")
			case strings.Contains(ex, "email"):
				return errz.NewBadRequest("email already exists")
			default:
				return errz.NewBadRequest("user already exists")
			}
		}
		return err
	}
	return nil
}

func (s RegistryService) RegisterCollegeAccount(req views.CollegeSignup, user *model.User) error {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Errorf("Failed to hash password: %v", err)
		return errz.NewBadRequest("failed to save user password")
	}
	var role model.Role
	err = s.db.Model(&model.Role{}).Where("name = ?", roles.College).First(&role).Error
	if err != nil {
		log.Errorf("Failed to find college role: %v", err)
		return errz.NewBadRequest("role not found")
	}
	var college model.College
	if err = s.db.Model(&model.College{}).Where("code = ?", req.Code).First(&college).Error; err != nil {
		return errz.NewBadRequest("college code not found")
	}

	var me = model.User{
		Username:     strings.TrimSpace(req.Username),
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: passwordHash,
		Roles:        []model.Role{role},
		CollegeID:    &college.ID,
		CreatedByID:  user.ID,
	}
	cleanedPhone := strings.TrimSpace(req.Phone)
	if cleanedPhone != "" {
		me.Phone = &cleanedPhone
	}
	err = db.DB.Model(&model.User{}).Create(&me).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ex := pgErr.Detail
			switch {
			case strings.Contains(ex, "username"):
				return errz.NewBadRequest("username already exists")
			case strings.Contains(ex, "email"):
				return errz.NewBadRequest("email already exists")
			default:
				return errz.NewBadRequest("user already exists")
			}
		}
		return err
	}
	return nil
}
