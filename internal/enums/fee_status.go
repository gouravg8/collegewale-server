package enums

import "fmt"

type FeeStatus string

const (
	FeeDue     FeeStatus = "due"
	FeePartial FeeStatus = "partial"
	FeePaid    FeeStatus = "paid"
	FeeOverdue FeeStatus = "overdue"
)

func (s FeeStatus) IsValid() error {
	switch s {
	case FeeDue, FeePartial, FeePaid, FeeOverdue:
	default:
		return fmt.Errorf("invalid fee status selected")
	}
	return nil
}
