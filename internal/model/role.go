package model

import (
	"collegeWaleServer/internal/enums/roles"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name roles.Roles `gorm:"type:varchar(80);not null;unique"`
}
