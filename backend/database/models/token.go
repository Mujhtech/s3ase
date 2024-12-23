package models

import (
	"time"

	"github.com/guregu/null"
)

type Token struct {
	ID        string      `json:"id" db:"id"`
	Value     string      `json:"value"  db:"value"`
	IsApp     bool        `json:"is_app" db:"is_app"`
	Type      string      `json:"type"  db:"type"`
	ExpiredAt *int64      `json:"expired_at" db:"expired_at"`
	IssuedAt  int64       `json:"issued_at" db:"issued_at"`
	Metadata  interface{} `json:"metadata" db:"metadata"`
	CreatedAt time.Time   `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at,omitempty"`
	DeletedAt null.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
