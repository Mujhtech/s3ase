package sse

import "encoding/json"

type EventType string

const (
	EventTypeUploadStarted   EventType = "upload_started"
	EventTypeUploadProgress  EventType = "upload_progress"
	EventTypeUploadCompleted EventType = "upload_completed"
	EventTypeUploadFailed    EventType = "upload_failed"
	EventTypeUploadCancelled EventType = "upload_cancelled"
	EventTypeUploadDeleted   EventType = "upload_deleted"
)

type UploadProgressStatus string

const (
	UploadProgressStatusStarted   UploadProgressStatus = "started"
	UploadProgressStatusPending   UploadProgressStatus = "pending"
	UploadProgressStatusUploading UploadProgressStatus = "uploading"
	UploadProgressStatusFailed    UploadProgressStatus = "failed"
	UploadProgressStatusCompleted UploadProgressStatus = "completed"
	UploadProgressStatusCancelled UploadProgressStatus = "cancelled"
)

type UploadProgress struct {
	FileID   string               `json:"id"`
	Name     string               `json:"name"`
	FolderID string               `json:"folder_id,omitempty"`
	Progress int64                `json:"progress"`
	Status   UploadProgressStatus `json:"status,omitempty"`
}

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}
