package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/utils"
	"collegeWaleServer/internal/views"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u UserService) MyInfo(user *model.User) (*views.MyInfo, error) {
	if user == nil {
		return nil, errz.NewNotFound("user not found")
	}
	myInfo := views.NewMyInfo(*user)
	return &myInfo, nil
}

func (u UserService) UpdateUser(user *model.User, req views.UpdateUserRequest) (*views.MyInfo, error) {
	if user == nil {
		return nil, errz.NewNotFound("user not found")
	}
	// validate request using view validation
	if err := req.IsValid(); err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if req.Email != nil {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if req.Username != nil {
		updates["username"] = strings.TrimSpace(*req.Username)
	}
	if req.Phone != nil {
		p := strings.TrimSpace(*req.Phone)
		updates["phone"] = p
	}
	if req.Password != nil {
		hashed, err := utils.HashPassword(strings.TrimSpace(*req.Password))
		if err != nil {
			return nil, errz.NewBadRequest("failed to hash password")
		}
		updates["password_hash"] = hashed
	}

	tx := u.db.Begin()
	if len(updates) > 0 {
		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Handle college change by code (if provided)
	if req.CollegeCode != nil {
		code := strings.TrimSpace(*req.CollegeCode)
		var college model.College
		if err := tx.Where("code = ?", code).First(&college).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, errz.NewNotFound("college not found")
			}
			return nil, err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Update("college_id", college.ID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Handle roles replacement (if provided)
	if req.Roles != nil {
		var rolesToSet []model.Role
		for _, r := range *req.Roles {
			var role model.Role
			if err := tx.Where("name = ?", r).First(&role).Error; err != nil {
				tx.Rollback()
				if err == gorm.ErrRecordNotFound {
					return nil, errz.NewNotFound("role not found")
				}
				return nil, err
			}
			rolesToSet = append(rolesToSet, role)
		}
		// Use Association.Replace to set new roles
		userModel := model.User{Model: gorm.Model{ID: user.ID}}
		if err := tx.Model(&userModel).Association("Roles").Replace(&rolesToSet); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load updated user with associations
	var updated model.User
	if err := u.db.Preload("Roles").Preload("Student").Preload("College").First(&updated, user.ID).Error; err != nil {
		return nil, err
	}

	myInfo := views.NewMyInfo(updated)
	return &myInfo, nil
}
