package service

import (
	"collegeWaleServer/internal/enums"
	"collegeWaleServer/internal/models"
	"collegeWaleServer/internal/services/email"
	"collegeWaleServer/internal/utils"
	auth_view "collegeWaleServer/internal/views/auth"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) CollegeSignup(req auth_view.CollegeSignup) (models.College, string, error) {
	// --- Input Validation ---
	if strings.TrimSpace(req.Name) == "" {
		return models.College{}, "", fmt.Errorf("college name cannot be empty")
	}
	if strings.TrimSpace(req.Email) == "" {
		return models.College{}, "", fmt.Errorf("email cannot be empty")
	}
	if !utils.IsEmailValid(req.Email) {
		return models.College{}, "", fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return models.College{}, "", fmt.Errorf("phone cannot be empty")
	}
	if !utils.IsPhoneValid(req.Phone) {
		return models.College{}, "", fmt.Errorf("invalid phone format")
	}
	if strings.TrimSpace(req.Code) == "" {
		return models.College{}, "", fmt.Errorf("college code cannot be empty")
	}
	if string(req.CourseType) == "" {
		return models.College{}, "", fmt.Errorf("course type cannot be empty")
	}
	switch req.CourseType {
	case enums.GNM, enums.ANM, enums.BSCNursing:
		// valid
	default:
		return models.College{}, "", fmt.Errorf("invalid course type: %s", req.CourseType)
	}
	if req.Seats <= 0 {
		return models.College{}, "", fmt.Errorf("seats must be greater than zero")
	}

	// --- gen token ---
	inviteToken := uuid.NewString()
	inviteTokenExpiry := time.Now().Add(24 * time.Hour)

	college := models.College{
		Name:         req.Name,
		Code:         req.Code,
		Phone:        req.Phone,
		Email:        req.Email,
		CourseType:   req.CourseType,
		Seats:        req.Seats,
		InviteToken:  inviteToken,
		InviteExpiry: inviteTokenExpiry,
		Logo:         req.Logo,
	}

	var existing models.College
	err := s.DB.Where("code = ?", req.Code).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// New college → create
			if err := s.DB.Create(&college).Error; err != nil {
				if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
					return models.College{}, "", fmt.Errorf("college with this name/email/code/phone already exists")
				}
				log.Error("Error creating college", err)
				return models.College{}, "", err
			}
		} else {
			return models.College{}, "", err
		}
	} else {
		// if clg exists → update token
		if err := s.DB.Model(&models.College{}).
			Where("code = ?", existing.Code).
			Updates(map[string]any{
				"invite_token":  inviteToken,
				"invite_expiry": inviteTokenExpiry,
			}).Error; err != nil {
			return models.College{}, "", err
		}
	}

	// --- send email ---
	emailService := email.NewEmailService()
	baseUrl := os.Getenv("APP_BASE_URL")
	data := map[string]string{
		"Name":             req.Name,
		"VerificationLink": baseUrl + "/verification?token=" + inviteToken,
	}

	if err := emailService.SendTemplateEmail(
		req.Email,
		"Verify Your College Account",
		"internal/services/email/templates/verification.html",
		data,
	); err != nil {
		return college, "", err
	}

	return college, "Verification email has been sent successfully", nil
}

func (s *AuthService) GetCollegeByToken(token string) (models.College, error) {
	if token == "" {
		return models.College{}, fmt.Errorf("Token is required")
	}

	var alreadyCollegeByToken models.College
	if err := s.DB.Where("token = ?", token).First(&alreadyCollegeByToken).Updates(map[string]any{
		"invite_token":  "",
		"invite_expiry": "",
	}).Error; err != nil {
		return models.College{}, err
	}

	return alreadyCollegeByToken, nil
}

func (s *AuthService) SetPassword(req auth_view.SetPassword) error {
	// todo: perform the hasing password
	var passwordHash string = ""

	err := s.DB.Where("code = ?", req.CollegeID).Updates(map[string]any{
		"password_hash": passwordHash,
	})

	if err != nil {
		return err.Error
	}
	return nil
}
