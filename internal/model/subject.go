package model

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name    string `gorm:"size:100;uniqueIndex;not null"`
	Code    string `gorm:"size:50;uniqueIndex;not null"`
	Credits int
}
