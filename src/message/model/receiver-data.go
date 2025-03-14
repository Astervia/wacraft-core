package message_model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	msg_model "github.com/Rfluid/whatsapp-cloud-api/src/message/model"
)

type ReceiverData struct {
	*msg_model.MessageReceived
}

// Implement the sql.Scanner interface
func (sd *ReceiverData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ProductData: expected []byte, got %T", value)
	}

	// Unmarshal the JSON-encoded data into the struct
	return json.Unmarshal(bytes, sd)
}

// Implement the driver.Valuer interface
func (sd ReceiverData) Value() (driver.Value, error) {
	// Marshal the struct into JSON
	if sd.MessageReceived == nil {
		return nil, nil
	}
	return json.Marshal(sd)
}
