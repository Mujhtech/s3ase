package models

import (
	"time"

	"github.com/guregu/null"
)

type FileStatus string

const (
	FileStatusPending   FileStatus = "pending"
	FileStatusUploading FileStatus = "uploading"
	FileStatusCompleted FileStatus = "completed"
	FileStatusFailed    FileStatus = "failed"
	FileStatusCancelled FileStatus = "cancelled"
	FileStatusStarted   FileStatus = "started"
)

type File struct {
	ID         string     `json:"id" db:"id"`
	FolderID   string     `json:"folder_id" db:"folder_id"`
	AppID      string     `json:"app_id" db:"app_id"`
	UploadedBy string     `json:"uploaded_by" db:"uploaded_by"`
	Name       string     `json:"name" db:"name"`
	Size       int64      `json:"size" db:"size"`
	Extension  string     `json:"extension" db:"extension"`
	MimeType   string     `json:"mime_type" db:"mime_type"`
	IsPublic   bool       `json:"is_public" db:"is_public"`
	PublicID   string     `json:"public_id" db:"public_id"`
	Status     FileStatus `json:"status" db:"status"`

	Metadata  interface{} `json:"metadata" db:"metadata"`
	CreatedAt time.Time   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
