package college

import (
	"collegeWaleServer/errz"
	"fmt"
)

type CourseType string

const (
	GNM        CourseType = "gnm"
	ANM        CourseType = "anm"
	BSCNursing CourseType = "bsc_nursing"
)

func (c CourseType) IsValidCourseType() error {
	switch c {
	case GNM, ANM, BSCNursing:
		return nil
	default:
		return errz.NewBadRequest(fmt.Sprintf("course type %s is not supported", c))
	}
}
