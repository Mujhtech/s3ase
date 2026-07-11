package dto

import "github.com/mujhtech/s3ase/database/models"

type CreateMemberRequestDto struct {
	Email string               `json:"email"`
	Role  models.AppMemberRole `json:"role"`
}
