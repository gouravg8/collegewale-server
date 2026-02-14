package utils

import (
	"regexp"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func IsPhoneValid(p string) bool {
	return phoneRegex.MatchString(p)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Warn("error in generating password", err)
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
