package enums

import (
	"fmt"
)

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

func (s StudentStatus) IsAllowed(ss StudentStatus) bool {
	switch s {
	case Draft:
		return ss == Submitted || ss == Verified || ss == Approved || ss == Unpayed || ss == Admitted
	case Submitted:
		return ss == Verified || ss == Approved || ss == Unpayed || ss == Admitted
	case Verified:
		return ss == Approved || ss == Unpayed || ss == Admitted
	case Approved:
		return ss == Unpayed || ss == Admitted
	case Unpayed:
		return ss == Admitted
	default:
		return false
	}
}

func (s StudentStatus) GetAllowedStatus() []StudentStatus {
	switch s {
	case Submitted:
		return []StudentStatus{Draft}
	case Verified:
		return []StudentStatus{Draft, Submitted}
	case Approved:
		return []StudentStatus{Draft, Submitted, Verified}
	case Unpayed:
		return []StudentStatus{Draft, Submitted, Verified, Approved}
	default:
		return []StudentStatus{Draft, Submitted, Verified, Approved} //for Admitted
	}
}
