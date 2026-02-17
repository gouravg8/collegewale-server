package common

type UserI interface {
	GetUsername() string
	GetPassword() string
	GetPhone() string
	GetEmail() string
}
