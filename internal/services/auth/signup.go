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
	if err := validateCollegeSignupRequest(req); err != nil {
		return models.College{}, "", err
	}

	// Validate APP_BASE_URL is configured
	baseUrl := os.Getenv("APP_BASE_URL")
	if baseUrl == "" {
		log.Error("APP_BASE_URL not configured")
		return models.College{}, "", fmt.Errorf("server configuration error: missing APP_BASE_URL")
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
				return models.College{}, "", fmt.Errorf("failed to create college: %w", err)
			}
		} else {
			log.Error("Database error checking existing college", err)
			return models.College{}, "", fmt.Errorf("signup failed: %w", err)
		}
	} else {
		// if clg exists → update token
		if err := s.DB.Model(&existing).
			Updates(map[string]any{
				"invite_token":  inviteToken,
				"invite_expiry": inviteTokenExpiry,
			}).Error; err != nil {
			log.Error("Error updating college invite token", err)
			return models.College{}, "", fmt.Errorf("failed to update invitation: %w", err)
		}
	}

	// --- send email ---
	emailService := email.NewEmailService()
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
		log.Error("Error sending verification email", err)
		return college, "", fmt.Errorf("failed to send verification email: %w", err)
	}

	return college, "Verification email has been sent successfully", nil
}

func (s *AuthService) GetCollegeByToken(token string) (models.College, error) {
	if token == "" {
		return models.College{}, fmt.Errorf("token cannot be empty")
	}

	var college models.College
	err := s.DB.Where("invite_token = ?", token).First(&college).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.College{}, fmt.Errorf("invalid token")
		}
		log.Error("Database error fetching college by token", err)
		return models.College{}, fmt.Errorf("failed to verify token: %w", err)
	}

	if college.InviteToken == "" {
		return models.College{}, fmt.Errorf("token has already been used")
	}

	if college.InviteExpiry.Before(time.Now()) {
		return models.College{}, fmt.Errorf("token has expired")
	}

	// Clear the token after verification
	if err := s.DB.Model(&college).Updates(map[string]any{
		"invite_token":  "",
		"invite_expiry": time.Time{},
	}).Error; err != nil {
		log.Error("Error clearing invite token", err)
		return models.College{}, fmt.Errorf("failed to complete verification: %w", err)
	}

	return college, nil
}

func (s *AuthService) SetPassword(req auth_view.SetPassword) error {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error("Error hashing password", err)
		return fmt.Errorf("failed to process password: %w", err)
	}

	var updateErr error
	if req.Code != "" {
		updateErr = s.DB.Where("code = ?", req.Code).Model(&models.College{}).Update("password_hash", passwordHash).Error
	} else if req.Email != "" {
		updateErr = s.DB.Where("email = ?", req.Email).Model(&models.College{}).Update("password_hash", passwordHash).Error
	}

	if updateErr != nil {
		log.Error("Error updating password", updateErr)
		return fmt.Errorf("failed to update password: %w", updateErr)
	}

	return nil
}

func (s *AuthService) CollegeLogin(req auth_view.CollegeLogin) (*models.College, error) {
	var college models.College
	var err error

	if req.Code != "" {
		err = s.DB.Where("code = ?", req.Code).First(&college).Error
	} else if req.Email != "" {
		err = s.DB.Where("email = ?", req.Email).First(&college).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Use timing attack resistant verification to avoid leaking user existence
			utils.VerifyPassword("$2a$10$dummyhashtoavoidtimingattack", req.Password)
			return nil, fmt.Errorf("invalid credentials")
		}
		log.Error("Database error during login", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if college.PasswordHash == "" {
		// College exists but password not set
		return nil, fmt.Errorf("password not set for this college - please complete verification")
	}

	if verifyErr := utils.VerifyPassword(college.PasswordHash, req.Password); verifyErr != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &college, nil
}

// validateCollegeSignupRequest validates the college signup request
func validateCollegeSignupRequest(req auth_view.CollegeSignup) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("college name cannot be empty")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !utils.IsEmailValid(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	if strings.TrimSpace(req.Phone) == "" {
		return fmt.Errorf("phone cannot be empty")
	}

	if !utils.IsPhoneValid(req.Phone) {
		return fmt.Errorf("invalid phone format")
	}

	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("college code cannot be empty")
	}

	if string(req.CourseType) == "" {
		return fmt.Errorf("course type cannot be empty")
	}

	switch req.CourseType {
	case enums.GNM, enums.ANM, enums.BSCNursing:
		// valid
	default:
		return fmt.Errorf("invalid course type: %s", req.CourseType)
	}

	if req.Seats <= 0 {
		return fmt.Errorf("seats must be greater than zero")
	}

	return nil
}
