package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	// The cost factor  10 is a good balance between security and performance.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
