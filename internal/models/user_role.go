package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	UserId uint
	RoleId uint
}
