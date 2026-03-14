package webhook_entity

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	webhook_model "github.com/Astervia/wacraft-core/src/webhook/model"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

// CircuitState represents the state of a circuit breaker
type CircuitState string

const (
	CircuitClosed   CircuitState = "closed"
	CircuitOpen     CircuitState = "open"
	CircuitHalfOpen CircuitState = "half_open"
)

type Webhook struct {
	Url           string `json:"url,omitempty" gorm:"not null"`
	Authorization string `json:"authorization,omitempty" gorm:"default:null"`
	HttpMethod    string `json:"http_method,omitempty" gorm:"not null"`
	Timeout       *int   `json:"timeout,omitempty" gorm:"default:1"` // The timeout in seconds. 0 means no timeout

	Event       webhook_model.Event `json:"event,omitempty" gorm:"not null"`
	WorkspaceID *uuid.UUID          `json:"workspace_id,omitempty" gorm:"type:uuid;index"`

	Workspace *workspace_entity.Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Security
	SigningSecret  string `json:"-" gorm:"default:null"` // Never exposed in JSON
	SigningEnabled bool   `json:"signing_enabled,omitempty" gorm:"default:false"`

	// Reliability
	MaxRetries   *int `json:"max_retries,omitempty" gorm:"default:3"`
	RetryDelayMs *int `json:"retry_delay_ms,omitempty" gorm:"default:1000"`
	IsActive     *bool `json:"is_active,omitempty" gorm:"default:true"`

	// Custom headers
	CustomHeaders map[string]string `json:"custom_headers,omitempty" gorm:"serializer:json;type:jsonb"`

	// Event filtering
	EventFilter *webhook_model.EventFilter `json:"event_filter,omitempty" gorm:"serializer:json;type:jsonb"`

	// Circuit breaker
	CircuitState    CircuitState `json:"circuit_state,omitempty" gorm:"default:'closed'"`
	FailureCount    int          `json:"failure_count,omitempty" gorm:"default:0"`
	LastFailureAt   *time.Time   `json:"last_failure_at,omitempty"`
	CircuitOpenedAt *time.Time   `json:"circuit_opened_at,omitempty"`

	common_model.Audit
}

func (w *Webhook) NewRequest(payload any) (*http.Request, error) {
	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create a new request
	req, err := http.NewRequest(w.HttpMethod, w.Url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// Add the Authorization header, if it exists
	if w.Authorization != "" {
		req.Header.Add("Authorization", w.Authorization)
	}

	// Set the Content-Type header to application/json
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (w *Webhook) ExecuteRequest(payload any, client *http.Client) (WebhookLog, error) {
	// Create a new WebhookLog
	webhookLog := WebhookLog{
		Payload:   payload,
		WebhookID: w.ID,
		Webhook:   w,
	}

	// Create a new request
	req, err := w.NewRequest(payload)
	if err != nil {
		return webhookLog, err
	}

	// If a timeout is specified, wrap the request with a context that cancels after the timeout
	if w.Timeout != nil && *w.Timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), time.Duration(*w.Timeout)*time.Second)
		defer cancel() // Ensure resources are released when the context is done
		req = req.WithContext(ctx)
	}

	// Execute the request using the provided client
	resp, err := client.Do(req)
	if err != nil {
		return webhookLog, err
	}
	defer resp.Body.Close()

	// Read the response data
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return webhookLog, err
	}

	// Initialize a variable to hold the parsed JSON response or null if it fails
	var parsedResponse any
	err = json.Unmarshal(responseData, &parsedResponse)
	if err != nil {
		// If the response is not valid JSON, set parsedResponse to nil (null)
		parsedResponse = nil
	}

	// Attach the response data to the webhookLog
	webhookLog.HttpResponseCode = resp.StatusCode
	webhookLog.ResponseData = parsedResponse

	return webhookLog, nil
}
