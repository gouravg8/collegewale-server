package views

type StudentSignUp struct { //TODO add more info as required per user
	Username *string `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}
