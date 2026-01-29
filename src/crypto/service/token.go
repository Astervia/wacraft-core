package crypto_service

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a token for email verification
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32)
}

// GenerateInvitationToken generates a token for workspace invitations
func GenerateInvitationToken() (string, error) {
	return GenerateSecureToken(32)
}

// GeneratePasswordResetToken generates a token for password reset
func GeneratePasswordResetToken() (string, error) {
	return GenerateSecureToken(32)
}
