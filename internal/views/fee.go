package views

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"
	"strings"
	"time"
)

type CreateFeeRecordRequest struct {
	StudentID   uint    `json:"student_id"`
	Term        string  `json:"term"`
	TotalAmount float64 `json:"total_amount"`
	DueDate     string  `json:"due_date"` // YYYY-MM-DD
}

func (r *CreateFeeRecordRequest) IsValid() error {
	r.Term = strings.TrimSpace(r.Term)
	r.DueDate = strings.TrimSpace(r.DueDate)

	if r.StudentID == 0 {
		return errz.NewBadRequest("student is required")
	}
	if r.Term == "" {
		return errz.NewBadRequest("term is required")
	}
	if r.TotalAmount <= 0 {
		return errz.NewBadRequest("total amount must be greater than zero")
	}
	if r.DueDate == "" {
		return errz.NewBadRequest("due date is required")
	}
	if _, err := time.Parse("2006-01-02", r.DueDate); err != nil {
		return errz.NewBadRequest("due date must be in YYYY-MM-DD format")
	}
	return nil
}

type RecordPaymentRequest struct {
	FeeRecordID uint    `json:"fee_record_id"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Note        string  `json:"note"`
}

func (r *RecordPaymentRequest) IsValid() error {
	r.Method = strings.TrimSpace(r.Method)
	r.Note = strings.TrimSpace(r.Note)

	if r.FeeRecordID == 0 {
		return errz.NewBadRequest("fee record is required")
	}
	if r.Amount <= 0 {
		return errz.NewBadRequest("payment amount must be greater than zero")
	}
	return nil
}

type PaymentResponse struct {
	ID     uint    `json:"id"`
	Amount float64 `json:"amount"`
	PaidAt string  `json:"paid_at"`
	Method string  `json:"method,omitempty"`
	Note   string  `json:"note,omitempty"`
}

type FeeRecordResponse struct {
	ID          uint              `json:"id"`
	StudentID   uint              `json:"student_id"`
	StudentName string            `json:"student_name"`
	RollNumber  string            `json:"roll_number"`
	Term        string            `json:"term"`
	TotalAmount float64           `json:"total_amount"`
	PaidAmount  float64           `json:"paid_amount"`
	Balance     float64           `json:"balance"`
	DueDate     string            `json:"due_date"`
	Status      string            `json:"status"`
	Payments    []PaymentResponse `json:"payments,omitempty"`
}

func NewFeeRecordResponse(f model.FeeRecord) FeeRecordResponse {
	paid := 0.0
	payments := make([]PaymentResponse, 0, len(f.Payments))
	for _, p := range f.Payments {
		paid += p.Amount
		payments = append(payments, PaymentResponse{
			ID:     p.ID,
			Amount: p.Amount,
			PaidAt: p.PaidAt.Format("2006-01-02"),
			Method: p.Method,
			Note:   p.Note,
		})
	}

	studentName := ""
	rollNumber := ""
	if f.Student.ID != 0 {
		studentName = strings.TrimSpace(f.Student.FirstName + " " + f.Student.LastName)
		rollNumber = f.Student.RollNumber
	}

	return FeeRecordResponse{
		ID:          f.ID,
		StudentID:   f.StudentID,
		StudentName: studentName,
		RollNumber:  rollNumber,
		Term:        f.Term,
		TotalAmount: f.TotalAmount,
		PaidAmount:  paid,
		Balance:     f.TotalAmount - paid,
		DueDate:     f.DueDate.Format("2006-01-02"),
		Status:      string(f.Status),
		Payments:    payments,
	}
}
