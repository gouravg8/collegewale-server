package enums

type StudentStatus string

const (
	Pending  StudentStatus = "pending"
	Approved StudentStatus = "approved"
	Unpayed  StudentStatus = "unpayed"
	Rejected StudentStatus = "rejected"
)
