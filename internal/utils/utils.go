package utils

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func IsPhoneValid(p string) bool {
	return phoneRegex.MatchString(p)
}
