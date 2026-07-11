package models

import (
	"time"

	"github.com/guregu/null"
)

type AppMemberRole string

const (
	AppMemberRoleOwner  AppMemberRole = "owner"
	AppMemberRoleMember AppMemberRole = "member"
)

type AppMember struct {
	ID        string        `json:"id" db:"id"`
	UserId    string        `json:"user_id" db:"user_id"`
	AppID     string        `json:"app_id" db:"app_id"`
	Role      AppMemberRole `json:"role" db:"role"`
	User      *User         `json:"user,omitempty"`
	Metadata  Metadata      `json:"metadata" db:"metadata"`
	CreatedAt time.Time     `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedAt null.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}
