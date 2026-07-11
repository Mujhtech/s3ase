package models

import (
	"time"

	"github.com/guregu/null"
)

type Webhook struct {
	ID          string      `json:"id" db:"id"`
	CreatedBy   string      `json:"created_by" db:"created_by"`
	AppID       string      `json:"app_id" db:"app_id"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
	URL         string      `json:"url" db:"url"`
	Metadata    Metadata    `json:"metadata" db:"metadata"`
	CreatedAt   time.Time   `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt   null.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
