package application

import (
	"collegeWaleServer/errz"
	"fmt"
)

type ApplicationStatus string

const (
	Draft     ApplicationStatus = "draft"
	Submitted ApplicationStatus = "submitted"
	Verified  ApplicationStatus = "verified"
	Approved  ApplicationStatus = "approved"
	Admitted  ApplicationStatus = "admitted"
	Rejected  ApplicationStatus = "rejected"
)

// allowed defines the state machine: current status → set of valid next statuses.
var allowed = map[ApplicationStatus][]ApplicationStatus{
	Draft:     {Submitted},
	Submitted: {Verified, Rejected},
	Verified:  {Approved, Rejected},
	Approved:  {Admitted, Rejected},
	Rejected:  {Draft},
	// Admitted is terminal — no transitions out.
}

func (s ApplicationStatus) IsValid() error {
	switch s {
	case Draft, Submitted, Verified, Approved, Admitted, Rejected:
		return nil
	default:
		return errz.NewBadRequest(fmt.Sprintf("invalid application status: %s", s))
	}
}

func (s ApplicationStatus) CanTransitionTo(next ApplicationStatus) error {
	if err := next.IsValid(); err != nil {
		return err
	}
	for _, a := range allowed[s] {
		if a == next {
			return nil
		}
	}
	return errz.NewBadRequest(fmt.Sprintf("cannot transition from '%s' to '%s'", s, next))
}
