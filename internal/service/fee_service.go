package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/enums"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/views"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FeeService struct {
	db *gorm.DB
}

func NewFeeService(db *gorm.DB) *FeeService {
	return &FeeService{db: db}
}

func deriveFeeStatus(totalAmount, paidAmount float64, dueDate time.Time) enums.FeeStatus {
	switch {
	case paidAmount >= totalAmount:
		return enums.FeePaid
	case paidAmount > 0:
		if time.Now().After(dueDate) {
			return enums.FeeOverdue
		}
		return enums.FeePartial
	case time.Now().After(dueDate):
		return enums.FeeOverdue
	default:
		return enums.FeeDue
	}
}

func (s *FeeService) CreateFeeRecord(req views.CreateFeeRecordRequest, collegeID uint, createdByID uint) (views.FeeRecordResponse, error) {
	var student model.Student
	if err := s.db.First(&student, req.StudentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.FeeRecordResponse{}, errz.NewNotFound("student not found")
		}
		return views.FeeRecordResponse{}, err
	}

	dueDate, _ := time.Parse("2006-01-02", req.DueDate)

	record := model.FeeRecord{
		StudentID:   req.StudentID,
		CollegeID:   collegeID,
		Term:        req.Term,
		TotalAmount: req.TotalAmount,
		DueDate:     dueDate,
		Status:      deriveFeeStatus(req.TotalAmount, 0, dueDate),
		CreatedByID: createdByID,
	}

	if err := s.db.Create(&record).Error; err != nil {
		return views.FeeRecordResponse{}, err
	}

	record.Student = student
	return views.NewFeeRecordResponse(record), nil
}

func (s *FeeService) ListFeeRecords(collegeID uint) ([]views.FeeRecordResponse, error) {
	var records []model.FeeRecord
	if err := s.db.
		Preload("Student").
		Preload("Payments").
		Where("college_id = ?", collegeID).
		Order("created_at desc").
		Find(&records).Error; err != nil {
		return nil, err
	}

	res := make([]views.FeeRecordResponse, 0, len(records))
	for _, r := range records {
		res = append(res, views.NewFeeRecordResponse(r))
	}
	return res, nil
}

func (s *FeeService) ListFeeRecordsForStudent(studentID uint) ([]views.FeeRecordResponse, error) {
	var records []model.FeeRecord
	if err := s.db.
		Preload("Student").
		Preload("Payments").
		Where("student_id = ?", studentID).
		Order("created_at desc").
		Find(&records).Error; err != nil {
		return nil, err
	}

	res := make([]views.FeeRecordResponse, 0, len(records))
	for _, r := range records {
		res = append(res, views.NewFeeRecordResponse(r))
	}
	return res, nil
}

func (s *FeeService) RecordPayment(req views.RecordPaymentRequest, collegeID uint, recordedByID uint) (views.FeeRecordResponse, error) {
	var updated model.FeeRecord

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var record model.FeeRecord
		if err := tx.Preload("Payments").Where("id = ? AND college_id = ?", req.FeeRecordID, collegeID).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errz.NewNotFound("fee record not found")
			}
			return err
		}

		payment := model.Payment{
			FeeRecordID:  record.ID,
			Amount:       req.Amount,
			PaidAt:       time.Now(),
			Method:       req.Method,
			Note:         req.Note,
			RecordedByID: recordedByID,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		record.Payments = append(record.Payments, payment)

		paidAmount := 0.0
		for _, p := range record.Payments {
			paidAmount += p.Amount
		}
		newStatus := deriveFeeStatus(record.TotalAmount, paidAmount, record.DueDate)
		if err := tx.Model(&record).Update("status", newStatus).Error; err != nil {
			return err
		}
		record.Status = newStatus

		if err := tx.Preload("Student").First(&record, record.ID).Error; err != nil {
			return err
		}
		updated = record
		return nil
	})

	if err != nil {
		return views.FeeRecordResponse{}, err
	}
	return views.NewFeeRecordResponse(updated), nil
}
