package views

type Me struct {
	Username *string `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

type Login struct {
	Role string
}
