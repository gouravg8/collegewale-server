package enums

type StudentStatus string

const (
	Draft     StudentStatus = "draft"
	Submitted StudentStatus = "submitted"
	Verified  StudentStatus = "verified"
	Approved  StudentStatus = "approved"
	Unpayed   StudentStatus = "unpayed"
	Admitted  StudentStatus = "admitted"
)
