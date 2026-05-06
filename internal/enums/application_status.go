package enums

import "fmt"

type StudentStatus string

const (
	Draft     StudentStatus = "draft"
	Submitted StudentStatus = "submitted"
	Verified  StudentStatus = "verified"
	Approved  StudentStatus = "approved"
	Unpayed   StudentStatus = "unpayed"
	Admitted  StudentStatus = "admitted"
)

func (s StudentStatus) IsValid() error {
	switch s {
	case Draft, Submitted, Verified, Approved, Unpayed, Admitted:
	default:
		return fmt.Errorf("invalid status selected")
	}
	return nil
}
