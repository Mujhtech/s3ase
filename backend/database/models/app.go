package models

import (
	"time"

	"github.com/guregu/null"
)

type App struct {
	ID          string      `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
	Slug        string      `json:"slug" db:"slug"`
	OwnerID     string      `json:"owner_id" db:"owner_id"`
	Region      null.String `json:"region" db:"region"`
	Bucket      string      `json:"bucket" db:"bucket"`
	Metadata    interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time   `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt   null.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
