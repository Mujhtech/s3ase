package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Metadata map[string]interface{}

func (m Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(m))
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	return decodeMetadata(data, m)
}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = Metadata{}
		return nil
	}
	var data []byte
	switch typed := value.(type) {
	case []byte:
		data = typed
	case string:
		data = []byte(typed)
	default:
		return fmt.Errorf("unsupported metadata value %T", value)
	}
	if len(data) == 0 {
		*m = Metadata{}
		return nil
	}
	return decodeMetadata(data, m)
}

func decodeMetadata(data []byte, destination *Metadata) error {
	var decoded interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	switch value := decoded.(type) {
	case map[string]interface{}:
		*destination = Metadata(value)
		return nil
	case string:
		// Older builds stored []byte JSON through database/sql, which PostgreSQL
		// encoded as a base64 JSON string (for example, "e30=" for {}).
		if legacy, err := base64.StdEncoding.DecodeString(value); err == nil {
			var object map[string]interface{}
			if err := json.Unmarshal(legacy, &object); err == nil {
				*destination = Metadata(object)
				return nil
			}
		}
		var object map[string]interface{}
		if err := json.Unmarshal([]byte(value), &object); err == nil {
			*destination = Metadata(object)
			return nil
		}
	}
	return fmt.Errorf("metadata must be a JSON object")
}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string]interface{}(m))
}
