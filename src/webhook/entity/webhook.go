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
)

type Webhook struct {
	Url           string `json:"url,omitempty" gorm:"not null"`
	Authorization string `json:"authorization,omitempty" gorm:"default:null"`
	HttpMethod    string `json:"http_method,omitempty" gorm:"not null"`
	Timeout       *int   `json:"timeout,omitempty" gorm:"default:1"` // The timeout in seconds. 0 means no timeout

	Event webhook_model.Event `json:"event,omitempty" gorm:"not null"`

	common_model.Audit
}

func (w *Webhook) NewRequest(payload interface{}) (*http.Request, error) {
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

func (w *Webhook) ExecuteRequest(payload interface{}, client *http.Client) (WebhookLog, error) {
	// Create a new WebhookLog
	webhookLog := WebhookLog{
		Payload:   payload,
		WebhookId: w.Id,
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
