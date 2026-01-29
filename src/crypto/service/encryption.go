package crypto_service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKeyLength   = errors.New("encryption key must be 32 bytes (256 bits)")
	ErrInvalidCiphertext  = errors.New("invalid ciphertext")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

// Encrypt encrypts plaintext using AES-256-GCM with the provided key.
// The key must be exactly 32 bytes.
// Returns base64-encoded ciphertext (nonce + encrypted data).
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM with the provided key.
// The key must be exactly 32 bytes.
func Decrypt(ciphertextB64 string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeyLength
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrCiphertextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateEncryptionKey generates a random 32-byte key for AES-256.
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// KeyFromBase64 decodes a base64-encoded key string to bytes.
func KeyFromBase64(keyB64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(keyB64)
}

// KeyToBase64 encodes a key to base64 string.
func KeyToBase64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
