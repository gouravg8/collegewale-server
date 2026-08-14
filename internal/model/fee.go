package model

import (
	"collegeWaleServer/internal/enums"
	"time"

	"gorm.io/gorm"
)

type FeeRecord struct {
	gorm.Model
	StudentID   uint
	Student     Student `gorm:"foreignKey:StudentID;references:ID;"`
	CollegeID   uint
	College     College `gorm:"foreignKey:CollegeID;references:ID;"`
	Term        string  `gorm:"size:80;not null"`
	TotalAmount float64 `gorm:"not null"`
	DueDate     time.Time
	Status      enums.FeeStatus `gorm:"size:20;not null;default:'due'"`
	CreatedByID uint
	Payments    []Payment `gorm:"foreignKey:FeeRecordID;"`
}

type Payment struct {
	gorm.Model
	FeeRecordID  uint
	Amount       float64   `gorm:"not null"`
	PaidAt       time.Time `gorm:"not null"`
	Method       string    `gorm:"size:40"`
	Note         string    `gorm:"size:255"`
	RecordedByID uint
}
