package message_model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	msg_model "github.com/Rfluid/whatsapp-cloud-api/src/message/model"
)

type SenderData struct {
	*msg_model.Message
}

// Implement the sql.Scanner interface
func (sd *SenderData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SenderData: expected []byte, got %T", value)
	}

	// Unmarshal the JSON-encoded data into the struct
	return json.Unmarshal(bytes, sd)
}

// Implement the driver.Valuer interface
func (sd SenderData) Value() (driver.Value, error) {
	// Marshal the struct into JSON
	if sd.Message == nil {
		return nil, nil
	}
	return json.Marshal(sd)
}
