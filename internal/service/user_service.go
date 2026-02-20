package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/views"

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
