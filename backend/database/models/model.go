package models

import "encoding/json"

type Metadata map[string]interface{}

func (m Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
