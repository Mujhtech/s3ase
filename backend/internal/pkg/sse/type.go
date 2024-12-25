package sse

import "encoding/json"

type EventType string

const (
	EventTypeUploadStarted  EventType = "upload_started"
	EventTypeUploadProgress EventType = "upload_progress"
	EventTypeUploadComplete EventType = "upload_complete"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}
