package models

import (
	"collegeWaleServer/internal/enums"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name enums.Roles `gorm:"type:varchar(80);not null;unique"`
}
