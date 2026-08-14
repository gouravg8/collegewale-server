package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/db"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/utils"
	"collegeWaleServer/internal/views"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var nonAlphaNumeric = regexp.MustCompile(`[^A-Z0-9]`)

func generateCollegeCode(name string) string {
	slug := nonAlphaNumeric.ReplaceAllString(strings.ToUpper(strings.TrimSpace(name)), "")
	if len(slug) > 6 {
		slug = slug[:6]
	}
	if slug == "" {
		slug = "CLG"
	}
	suffix := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))[:4]
	return slug + suffix
}

type RegistryService struct {
	db *gorm.DB
}

func NewRegistryService(db *gorm.DB) *RegistryService {
	return &RegistryService{db}
}

func (s RegistryService) RegisterCollege(req views.CollegeRequest, user *model.User, objectKey string) error {
	var courses []model.Courses
	if err := s.db.Model(&model.Courses{}).Find(&courses).Error; err != nil {
		return err
	}
	var courseList []model.Courses
	for _, rc := range req.Courses {
		found := false
		for _, c := range courses {
			if c.Name == rc {
				courseList = append(courseList, c)
				found = true
				break
			}
		}
		if !found {
			return errz.NewBadRequest(fmt.Sprintf("course :: %s not found", rc))
		}
	}
	if len(courseList) != len(req.Courses) {
		return errz.NewBadRequest("courses not found")
	}
	clg := model.College{
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Phone:       strings.TrimSpace(req.Phone),
		Email:       strings.TrimSpace(req.Email),
		Courses:     courseList,
		Seats:       req.Seats,
		Logo:        objectKey,
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
	var courseID uint
	err = s.db.Model(&model.Courses{}).Where("name = ?", strings.ToUpper(req.CourseType)).Pluck("id", &courseID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errz.NewBadRequest("course not found")
		}
		return err
	}
	var student = model.Student{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Phone:            req.Phone,
		RollNumber:       req.RollNumber,
		CourseID:         courseID,
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
	err = s.db.Model(&model.Role{}).Where("name = ?", roles.CollegeAdmin).First(&role).Error
	if err != nil {
		log.Errorf("Failed to find college_admin role: %v", err)
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

// RegisterCollegeWithAdmin creates a College with sane defaults (auto
// generated code, zero courses, a nominal seat count) and its College
// Admin login user in a single transaction -- the minimal flow where the
// platform Admin only supplies a college name and the admin's own details.
func (s RegistryService) RegisterCollegeWithAdmin(req views.CollegeWithAdminRequest, createdBy *model.User) error {
	passwordHash, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		log.Errorf("Failed to hash password: %v", err)
		return errz.NewBadRequest("failed to save admin password")
	}

	var role model.Role
	if err := s.db.Model(&model.Role{}).Where("name = ?", roles.CollegeAdmin).First(&role).Error; err != nil {
		log.Errorf("Failed to find college_admin role: %v", err)
		return errz.NewBadRequest("role not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		college := model.College{
			Name:        strings.TrimSpace(req.CollegeName),
			Code:        generateCollegeCode(req.CollegeName),
			Phone:       strings.TrimSpace(req.AdminPhone),
			Email:       strings.TrimSpace(req.AdminEmail),
			Seats:       1,
			CreatedById: createdBy.ID,
		}
		if err := tx.Create(&college).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return errz.NewBadRequest("a college with this name already exists")
			}
			return err
		}

		admin := model.User{
			Username:     strings.TrimSpace(req.AdminUsername),
			Email:        strings.TrimSpace(req.AdminEmail),
			PasswordHash: passwordHash,
			Roles:        []model.Role{role},
			CollegeID:    &college.ID,
			CreatedByID:  createdBy.ID,
		}
		cleanedPhone := strings.TrimSpace(req.AdminPhone)
		if cleanedPhone != "" {
			admin.Phone = &cleanedPhone
		}
		if err := tx.Create(&admin).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				ex := pgErr.Detail
				switch {
				case strings.Contains(ex, "username"):
					return errz.NewBadRequest("username already exists")
				case strings.Contains(ex, "email"):
					return errz.NewBadRequest("email already exists")
				default:
					return errz.NewBadRequest("admin user already exists")
				}
			}
			return err
		}
		return nil
	})
}
