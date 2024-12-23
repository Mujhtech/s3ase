package dto

import "github.com/mujhtech/s3ase/database/models"

type CreateApiKeyRequestDto struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	ExpiredAt   int64               `json:"expired_at"`
	Access      models.ApiKeyAccess `json:"access"`
}
