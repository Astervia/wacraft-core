package crypto_service

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(truncateTo72Bytes(password)), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Ensures the password is no longer than 72 bytes.
func truncateTo72Bytes(password string) string {
	if len(password) > 72 {
		return password[:72]
	}
	return password
}
