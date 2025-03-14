package status_model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	webhook_model "github.com/Rfluid/whatsapp-cloud-api/src/webhook/model"
)

type ProductData struct {
	*webhook_model.Status
}

// Implement the sql.Scanner interface
func (pd *ProductData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ProductData: expected []byte, got %T", value)
	}

	// Unmarshal the JSON-encoded data into the struct
	return json.Unmarshal(bytes, pd)
}

// Implement the driver.Valuer interface
func (pd ProductData) Value() (driver.Value, error) {
	// Marshal the struct into JSON
	return json.Marshal(pd)
}
