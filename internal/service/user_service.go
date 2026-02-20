package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u UserService) MyInfo(user *model.User) (interface{}, error) {
	if user == nil {
		return nil, errz.NewNotFound("user not found")
	}
	return user, nil
}
