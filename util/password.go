package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	if hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost); err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	} else {
		return string(hashed), nil
	}
}
func CheckPassword(pw, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
}
