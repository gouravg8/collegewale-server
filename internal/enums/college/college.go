package college

import (
	"collegeWaleServer/errz"
	"fmt"
)

type CourseType string

const (
	GNM        CourseType = "GNM"
	ANM        CourseType = "ANM"
	BSCNursing CourseType = "BSc Nursing"
)

func (c CourseType) IsValidCourseType() error {
	switch c {
	case GNM, ANM, BSCNursing:
		return nil
	default:
		return errz.NewBadRequest(fmt.Sprintf("course type %s is not supported", c))
	}
}
