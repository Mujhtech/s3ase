package models

import (
	"time"

	"github.com/guregu/null"
)

type AppSubscription struct {
	ID        string    `json:"id" db:"id"`
	AppID     string    `json:"app_id" db:"app_id"`
	Plan      string    `json:"plan" db:"plan"`
	Status    string    `json:"status" db:"status"`
	Metadata  Metadata  `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
