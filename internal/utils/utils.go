package utils

import (
	"encoding/base64"
	"regexp"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func IsPhoneValid(p string) bool {
	return phoneRegex.MatchString(p)
}

func HashPassword(password string) (string, error) {

	decoded, err := base64.StdEncoding.DecodeString(password)

	if err != nil {
		return "", err
	}

	bytes, err := bcrypt.GenerateFromPassword(decoded, 14)
	if err != nil {
		log.Warn("error in generating password", err)
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(password, hash string) bool {
	decoded, err := base64.StdEncoding.DecodeString(password)

	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), decoded)
	return err == nil
}
