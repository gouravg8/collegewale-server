package enums

type CourseType string

const (
	GNM        CourseType = "GNM"
	ANM        CourseType = "ANM"
	BSCNursing CourseType = "BSc Nursing"
)

type CollegeType string

const (
	PENDING  CollegeType = "PENDING"
	ACTIVE   CollegeType = "ACTIVE"
	REJECTED CollegeType = "REJECTED"
)
