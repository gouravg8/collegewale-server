package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/services/email"
	"collegeWaleServer/internal/utils"
	"collegeWaleServer/internal/views"
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
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) CollegeSignup(req auth_view.CollegeSignup) (model.College, string, error) {
	// --- Input Validation ---
	if strings.TrimSpace(req.Name) == "" {
		return model.College{}, "", fmt.Errorf("college name cannot be empty")
	}
	if strings.TrimSpace(req.Email) == "" {
		return model.College{}, "", fmt.Errorf("email cannot be empty")
	}
	if !utils.IsEmailValid(req.Email) {
		return model.College{}, "", fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return model.College{}, "", fmt.Errorf("phone cannot be empty")
	}
	if !utils.IsPhoneValid(req.Phone) {
		return model.College{}, "", fmt.Errorf("invalid phone format")
	}
	if strings.TrimSpace(req.Code) == "" {
		return model.College{}, "", fmt.Errorf("college code cannot be empty")
	}
	if string(req.CourseType) == "" {
		return model.College{}, "", fmt.Errorf("course type cannot be empty")
	}
	if err := req.CourseType.IsValidCourseType(); err != nil {
		return model.College{}, "", err
	}
	if req.Seats <= 0 {
		return model.College{}, "", fmt.Errorf("seats must be greater than zero")
	}

	// --- gen token ---
	inviteToken := uuid.NewString()
	inviteTokenExpiry := time.Now().Add(24 * time.Hour)

	college := model.College{
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

	var existing model.College
	err := s.db.Where("code = ?", req.Code).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// New college → create
			if err := s.db.Create(&college).Error; err != nil {
				if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
					return model.College{}, "", fmt.Errorf("college with this name/email/code/phone already exists")
				}
				log.Error("Error creating college", err)
				return model.College{}, "", err
			}
		} else {
			return model.College{}, "", err
		}
	} else {
		// if clg exists → update token
		if err := s.db.Model(&model.College{}).
			Where("code = ?", existing.Code).
			Updates(map[string]any{
				"invite_token":  inviteToken,
				"invite_expiry": inviteTokenExpiry,
			}).Error; err != nil {
			return model.College{}, "", err
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

func (s *AuthService) GetCollegeByToken(token string) (model.College, error) {
	if token == "" {
		return model.College{}, fmt.Errorf("Token is required")
	}

	var alreadyCollegeByToken model.College
	if err := s.db.Where("token = ?", token).First(&alreadyCollegeByToken).Updates(map[string]any{
		"invite_token":  "",
		"invite_expiry": "",
	}).Error; err != nil {
		return model.College{}, err
	}

	return alreadyCollegeByToken, nil
}

func (s *AuthService) SetPassword(req auth_view.SetPassword) error {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	if req.Code != "" {
		err = s.db.Where("code = ?", req.Code).Updates(map[string]any{
			"password_hash": passwordHash,
		}).Error
	} else if req.Email != "" {
		err = s.db.Where("email = ?", req.Email).Updates(map[string]any{
			"password_hash": passwordHash,
		}).Error
	}

	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) CollegeLogin(req auth_view.CollegeLogin) (*model.College, error) {
	var college model.College

	if req.Code != "" {
		if err := s.db.Where("code = ?", req.Code).First(&college).Error; err != nil {
			return &model.College{}, fmt.Errorf("Error %v", err.Error())
		}
	} else if req.Email != "" {
		if err := s.db.Where("email = ?", req.Email).First(&college).Error; err != nil {
			return &model.College{}, fmt.Errorf("Error %v", err.Error())
		}
	}

	return &college, nil
}

func (s *AuthService) SignIn(req views.MeLogin) (*views.Me, error) {
	var me model.User
	q := s.db.Model(&model.User{})
	if req.Username != nil && *req.Username != "" {
		q.Where("username = ?", *req.Username)
	} else if req.Email != nil && *req.Email != "" {
		q.Where("email = ?", req.Email)
	} else if req.Phone != nil && *req.Phone != "" {
		q.Where("phone = ?", *req.Phone)
	} else {
		return nil, errz.NewBadRequest("Please provide a valid Username, Email or Phone.")
	}
	if err := q.Preload("Role").First(&me).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errz.NewNotFound("User not found")
		}
		return nil, err
	}
	if !utils.ComparePassword(req.Password, me.PasswordHash) {
		return nil, errz.NewBadRequest("Incorrect Password.")

	}
	res := &views.Me{
		Email: me.Email,
		Phone: me.Phone,
	}
	if me.CollegeID != nil {
		res.CollegeID = *me.CollegeID
	}
	return res, nil
}

//func getExistingColleges(db *gorm.db) (map[string]*model.College, error) {
//
//}
