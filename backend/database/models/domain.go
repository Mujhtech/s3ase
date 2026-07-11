package models

import (
	"time"

	"github.com/guregu/null"
)

type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "pending"
	DomainStatusVerified DomainStatus = "verified"
	DomainStatusFailed   DomainStatus = "failed"
	DomainStatusReview   DomainStatus = "review"
)

type Domain struct {
	ID          string       `json:"id" db:"id"`
	CreatedBy   string       `json:"created_by" db:"created_by"`
	AppID       string       `json:"app_id" db:"app_id"`
	Domain      string       `json:"domain" db:"domain"`
	CnameRecord string       `json:"cname_record" db:"cname_record"`
	TxtRecord   string       `json:"txt_record" db:"txt_record"`
	Metadata    Metadata     `json:"metadata" db:"metadata"`
	Status      DomainStatus `json:"status" db:"status"`
	CreatedAt   time.Time    `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt   null.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
}
