package models

import (
	"time"

	"github.com/guregu/null"
)

type Folder struct {
	ID          string      `json:"id" db:"id"`
	ParentID    null.String `json:"parent_id" db:"parent_id"`
	AppID       string      `json:"app_id" db:"app_id"`
	CreatedBy   string      `json:"created_by" db:"created_by"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
	Metadata    Metadata    `json:"metadata" db:"metadata"`
	CreatedAt   time.Time   `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt   null.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
