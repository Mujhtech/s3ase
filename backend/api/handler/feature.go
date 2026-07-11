package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/pkg/response"
)

type Feature struct {
	Name                      string `json:"name"`
	Description               string `json:"description"`
	IsGithubAuthEnabled       bool   `json:"is_github_auth_enabled"`
	IsGoogleAuthEnabled       bool   `json:"is_google_auth_enabled"`
	IsAwsConfigured           bool   `json:"is_aws_configured"`
	IsObjectStorageConfigured bool   `json:"is_object_storage_configured"`
	ObjectStorageProvider     string `json:"object_storage_provider"`
	ObjectStorageRegion       string `json:"object_storage_region"`
	Version                   string `json:"version"`
}

func (h *Handler) GetFeatures(w http.ResponseWriter, r *http.Request) {

	isGoogleAuthEnabled := h.cfg.Auth.GoogleAuth.ClientID != "" && h.cfg.Auth.GoogleAuth.ClientSecret != ""
	isGithubAuthEnabled := h.cfg.Auth.GithubAuth.ClientID != "" && h.cfg.Auth.GithubAuth.ClientSecret != ""
	isAwsConfigured := h.cfg.Aws.AccessKey != "" && h.cfg.Aws.SecretKey != ""
	isObjectStorageConfigured := isAwsConfigured
	if h.cfg.ObjectStorage.Provider == config.ObjectStorageProviderR2 {
		isObjectStorageConfigured = h.cfg.R2.AccountID != "" && h.cfg.R2.AccessKeyID != "" && h.cfg.R2.SecretKey != ""
	}

	_ = response.Ok(w, r, "feature data retrieved successfully", Feature{
		Name:                      "s3ase",
		Description:               "Simplifying S3 usage through open source",
		IsGithubAuthEnabled:       isGithubAuthEnabled,
		IsGoogleAuthEnabled:       isGoogleAuthEnabled,
		IsAwsConfigured:           isAwsConfigured,
		IsObjectStorageConfigured: isObjectStorageConfigured,
		ObjectStorageProvider:     string(h.cfg.ObjectStorage.Provider),
		ObjectStorageRegion:       h.cfg.ObjectStorageRegion(),
		Version:                   "0.0.1",
	})
}
