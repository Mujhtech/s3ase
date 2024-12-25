package auth

import (
	"github.com/mujhtech/s3ase/database/models"
)

type UserSession struct {
	User     *models.User
	Metadata *TokenMetadata
}

type AppSession struct {
	App      *models.App
	Metadata *TokenMetadata
}

type TokenMetadata struct {
	Type     CredentialType
	Metadata interface{}
	TokenID  string
}
