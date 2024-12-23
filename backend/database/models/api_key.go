package models

import (
	"time"

	"github.com/guregu/null"
)

type ApiKeyAccess string

const (
	ApiKeyAccessNone  ApiKeyAccess = "none"
	ApiKeyAccessRead  ApiKeyAccess = "read"
	ApiKeyAccessWrite ApiKeyAccess = "write"
	ApiKeyAccessFull  ApiKeyAccess = "full"
)

type ApiKey struct {
	ID          string       `json:"id" db:"id"`
	CreatedBy   string       `json:"created_by" db:"created_by"`
	AppID       string       `json:"app_id" db:"app_id"`
	Name        string       `json:"name" db:"name"`
	Description null.String  `json:"description" db:"description"`
	Access      ApiKeyAccess `json:"access" db:"access"`
	ExpiredAt   null.Int     `json:"expired_at" db:"expired_at"`
	LastUsed    null.Time    `json:"last_used" db:"last_used"`
	Metadata    interface{}  `json:"metadata" db:"metadata"`
	CreatedAt   time.Time    `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt   null.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
}
