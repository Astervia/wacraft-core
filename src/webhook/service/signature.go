package webhook_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// SignatureVersion is the current signature algorithm version
const SignatureVersion = "v1"

// GenerateSignature creates an HMAC-SHA256 signature for webhook payloads
// Format: v1={hex-encoded-hmac-sha256}
func GenerateSignature(secret string, timestamp int64, payload []byte) string {
	// Construct the message: "v1:{timestamp}:{payload}"
	message := fmt.Sprintf("%s:%d:%s", SignatureVersion, timestamp, string(payload))

	// Create HMAC-SHA256 hash
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	hash := h.Sum(nil)

	// Return in format v1={hex-encoded-hash}
	return fmt.Sprintf("%s=%s", SignatureVersion, hex.EncodeToString(hash))
}

// SignatureHeaders generates the signature headers for a webhook request
// Returns X-Wacraft-Signature and X-Wacraft-Timestamp header values
func SignatureHeaders(secret string, payload []byte) (signature string, timestamp string) {
	ts := time.Now().Unix()
	sig := GenerateSignature(secret, ts, payload)
	return sig, strconv.FormatInt(ts, 10)
}

// VerifySignature verifies an HMAC-SHA256 signature (constant-time comparison)
// This is provided for documentation/testing purposes - consumers implement their own
func VerifySignature(secret string, timestamp int64, payload []byte, signature string) bool {
	expected := GenerateSignature(secret, timestamp, payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// IsTimestampValid checks if a timestamp is within the allowed time window
// Default window is 5 minutes (300 seconds)
func IsTimestampValid(timestamp int64, maxAgeSeconds int64) bool {
	if maxAgeSeconds <= 0 {
		maxAgeSeconds = 300 // 5 minutes default
	}
	now := time.Now().Unix()
	return now-timestamp <= maxAgeSeconds
}
