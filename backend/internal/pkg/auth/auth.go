package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mujhtech/s3ase/config"
	"golang.org/x/oauth2"
)

type Emails []Email

type Email struct {
	Email    string
	Verified bool
	Primary  bool
}

type User struct {
	AuthenticationMethod string
	Emails               Emails
	Metadata             *Claims
}

type Claims struct {
	Issuer        string                 `json:"iss,omitempty" structs:"iss,omitempty"`
	Subject       string                 `json:"sub,omitempty" structs:"sub,omitempty"`
	Name          string                 `json:"name" structs:"name,omitempty"`
	FamilyName    string                 `json:"family_name" structs:"family_name,omitempty"`
	GivenName     string                 `json:"given_name" structs:"given_name,omitempty"`
	MiddleName    string                 `json:"middle_name" structs:"middle_name,omitempty"`
	Username      string                 `json:"username" structs:"username,omitempty"`
	Email         string                 `json:"email" structs:"email,omitempty"`
	EmailVerified bool                   `json:"email_verified" structs:"email_verified,omitempty"`
	AvatarUrl     string                 `json:"avatar_url" structs:"avatar_url,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type AuthProvider interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Name() string
	GetOAuthToken(string) (*oauth2.Token, error)
	GetUser(*oauth2.Token) (*User, error)
}

func GetAuthProvider(cfg *config.Config, name string, redirectUrl string) (AuthProvider, error) {
	switch name {
	case "google":
		if cfg.Auth.GoogleAuth.ClientID == "" || cfg.Auth.GoogleAuth.ClientSecret == "" {
			return nil, fmt.Errorf("google authentication is not configured")
		}
		return NewGoogleProvider(cfg.Auth.GoogleAuth, redirectUrl)
	case "github":
		if cfg.Auth.GithubAuth.ClientID == "" || cfg.Auth.GithubAuth.ClientSecret == "" {
			return nil, fmt.Errorf("github authentication is not configured")
		}
		return NewGithubProvider(cfg.Auth.GithubAuth, redirectUrl)
	default:
		return nil, fmt.Errorf("auth provider %s not supported", name)
	}
}

func makeRequest(token *oauth2.Token, config *oauth2.Config, url string, dst interface{}) error {
	client := config.Client(context.Background(), token)
	client.Timeout = time.Second * 10
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	bodyBytes, readErr := io.ReadAll(res.Body)
	closeErr := res.Body.Close()
	if readErr != nil {
		return fmt.Errorf("failed to read authentication response: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close authentication response: %w", closeErr)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return errors.New(string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, dst); err != nil {
		return err
	}

	return nil
}
