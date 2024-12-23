package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/internal/pkg/response"
)

type Feature struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	IsGithubAuthEnabled bool   `json:"is_github_auth_enabled"`
	IsGoogleAuthEnabled bool   `json:"is_google_auth_enabled"`
	IsAwsConfigured     bool   `json:"is_aws_configured"`
	Version             string `json:"version"`
}

func (h *Handler) GetFeatures(w http.ResponseWriter, r *http.Request) {

	isGoogleAuthEnabled := h.cfg.Auth.GoogleAuth.ClientID != "" && h.cfg.Auth.GoogleAuth.ClientSecret != ""
	isGithubAuthEnabled := h.cfg.Auth.GithubAuth.ClientID != "" && h.cfg.Auth.GithubAuth.ClientSecret != ""
	isAwsConfigured := h.cfg.Aws.AccessKey != "" && h.cfg.Aws.SecretKey != ""

	_ = response.Ok(w, r, "feature data retrieved successfully", Feature{
		Name:                "s3ase",
		Description:         "Simplifying S3 usage through open source",
		IsGithubAuthEnabled: isGithubAuthEnabled,
		IsGoogleAuthEnabled: isGoogleAuthEnabled,
		IsAwsConfigured:     isAwsConfigured,
		Version:             "0.0.1",
	})
}
