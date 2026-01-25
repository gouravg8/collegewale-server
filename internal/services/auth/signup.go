package service

import (
	"collegeWaleServer/internal/enums"
	"collegeWaleServer/internal/models"
	"collegeWaleServer/internal/utils"
	auth_view "collegeWaleServer/internal/views/auth"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) CollegeSignup(req auth_view.CollegeSignup) (models.College, error) {
	if strings.TrimSpace(req.Name) == "" {
		return models.College{}, fmt.Errorf("Empty name")
	}

	if strings.TrimSpace(req.Email) == "" {
		return models.College{}, fmt.Errorf("Empty email")
	}

	if !utils.IsEmailValid(req.Email) {
		return models.College{}, fmt.Errorf("email not valid")
	}

	if strings.TrimSpace(req.Phone) == "" {
		return models.College{}, fmt.Errorf("Empty phone")
	}

	if strings.TrimSpace(req.Code) == "" {
		return models.College{}, fmt.Errorf("college code cannot be empty")
	}

	if string(req.CourseType) == "" {
		return models.College{}, fmt.Errorf("college type cannot be empty")
	}

	switch req.CourseType {
	case enums.GNM, enums.ANM, enums.BSCNursing:
		// valid
	default:
		return models.College{}, fmt.Errorf("invalid course type: %s", req.CourseType)
	}

	if req.Seats == 0 {
		return models.College{}, fmt.Errorf("seats must be greater than zero")
	}

	college := models.College{
		Name:       req.Name,
		Code:       req.Code,
		Phone:      req.Phone,
		Email:      req.Email,
		CourseType: req.CourseType,
		Seats:      req.Seats,
		Logo:       req.Logo,
	}

	// err := s.db.Create(&college).Error
	// err := s.db.Create(&college).Error
	err := s.DB.Create(&college).Error

	if err != nil {
		log.Error("Error creating college", err)
		return models.College{}, err
	}

	return college, nil
}
