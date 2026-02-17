package model

import "gorm.io/gorm"

type Courses struct { //TODO use list of []courses as it can be expandable
	gorm.Model
	Name        string
	Description string `gorm:"type:text;default:''"`
}
