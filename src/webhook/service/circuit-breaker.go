package webhook_service

import (
	"time"

	webhook_entity "github.com/Astervia/wacraft-core/src/webhook/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Circuit breaker configuration
const (
	// FailureThreshold is the number of consecutive failures before opening the circuit
	FailureThreshold = 5
	// RecoveryTimeout is the duration before attempting to close the circuit
	RecoveryTimeout = 30 * time.Second
)

// CircuitBreaker provides circuit breaker functionality for webhooks
type CircuitBreaker struct {
	db *gorm.DB
}

// NewCircuitBreaker creates a new circuit breaker instance
func NewCircuitBreaker(db *gorm.DB) *CircuitBreaker {
	return &CircuitBreaker{db: db}
}

// AllowRequest checks if a request is allowed for the given webhook
// Returns true if the circuit is closed or half-open, false if open
func (cb *CircuitBreaker) AllowRequest(webhookID uuid.UUID) (bool, error) {
	var webhook webhook_entity.Webhook
	if err := cb.db.Select("circuit_state", "circuit_opened_at").First(&webhook, "id = ?", webhookID).Error; err != nil {
		return false, err
	}

	switch webhook.CircuitState {
	case webhook_entity.CircuitClosed:
		return true, nil

	case webhook_entity.CircuitHalfOpen:
		return true, nil

	case webhook_entity.CircuitOpen:
		// Check if recovery timeout has passed
		if webhook.CircuitOpenedAt != nil && time.Since(*webhook.CircuitOpenedAt) >= RecoveryTimeout {
			// Transition to half-open
			if err := cb.transitionToHalfOpen(webhookID); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil

	default:
		// Unknown state, treat as closed
		return true, nil
	}
}

// RecordSuccess records a successful request and potentially closes the circuit
func (cb *CircuitBreaker) RecordSuccess(webhookID uuid.UUID) error {
	return cb.db.Model(&webhook_entity.Webhook{}).
		Where("id = ?", webhookID).
		Updates(map[string]any{
			"circuit_state":     webhook_entity.CircuitClosed,
			"failure_count":     0,
			"last_failure_at":   nil,
			"circuit_opened_at": nil,
		}).Error
}

// RecordFailure records a failed request and potentially opens the circuit
func (cb *CircuitBreaker) RecordFailure(webhookID uuid.UUID) error {
	var webhook webhook_entity.Webhook
	if err := cb.db.Select("circuit_state", "failure_count").First(&webhook, "id = ?", webhookID).Error; err != nil {
		return err
	}

	now := time.Now()
	newFailureCount := webhook.FailureCount + 1

	// If we're in half-open state, any failure opens the circuit again
	if webhook.CircuitState == webhook_entity.CircuitHalfOpen {
		return cb.db.Model(&webhook_entity.Webhook{}).
			Where("id = ?", webhookID).
			Updates(map[string]any{
				"circuit_state":     webhook_entity.CircuitOpen,
				"failure_count":     newFailureCount,
				"last_failure_at":   now,
				"circuit_opened_at": now,
			}).Error
	}

	// Check if we've reached the failure threshold
	if newFailureCount >= FailureThreshold {
		return cb.db.Model(&webhook_entity.Webhook{}).
			Where("id = ?", webhookID).
			Updates(map[string]any{
				"circuit_state":     webhook_entity.CircuitOpen,
				"failure_count":     newFailureCount,
				"last_failure_at":   now,
				"circuit_opened_at": now,
			}).Error
	}

	// Just increment the failure count
	return cb.db.Model(&webhook_entity.Webhook{}).
		Where("id = ?", webhookID).
		Updates(map[string]any{
			"failure_count":   newFailureCount,
			"last_failure_at": now,
		}).Error
}

// transitionToHalfOpen transitions the circuit to half-open state
func (cb *CircuitBreaker) transitionToHalfOpen(webhookID uuid.UUID) error {
	return cb.db.Model(&webhook_entity.Webhook{}).
		Where("id = ?", webhookID).
		Update("circuit_state", webhook_entity.CircuitHalfOpen).Error
}

// GetState returns the current circuit state for a webhook
func (cb *CircuitBreaker) GetState(webhookID uuid.UUID) (webhook_entity.CircuitState, error) {
	var webhook webhook_entity.Webhook
	if err := cb.db.Select("circuit_state").First(&webhook, "id = ?", webhookID).Error; err != nil {
		return "", err
	}
	return webhook.CircuitState, nil
}

// ResetCircuit resets the circuit breaker to closed state
func (cb *CircuitBreaker) ResetCircuit(webhookID uuid.UUID) error {
	return cb.RecordSuccess(webhookID)
}
