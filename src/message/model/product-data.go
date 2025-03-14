package message_model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	message_model "github.com/Rfluid/whatsapp-cloud-api/src/message/model"
)

type ProductData struct {
	*message_model.Response
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
	if pd.Response == nil {
		return nil, nil
	}
	return json.Marshal(pd)
}
