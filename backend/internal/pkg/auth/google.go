package auth

import (
	"context"

	"github.com/mujhtech/s3ase/config"
	"golang.org/x/oauth2"
)

const (
	GoogleUrl                = "https://accounts.google.com"
	GoogleApiBase            = "https://www.googleapis.com"
	GoogleUserAPI            = GoogleApiBase + "/oauth2/v3/userinfo"
	GoogleOauthAuthEndpoint  = GoogleUrl + "/o/oauth2/auth"
	GoogleOauthTokenEndpoint = GoogleUrl + "/o/oauth2/token"
)

type GoogleUser struct {
	ID            string `json:"sub"`
	Name          string `json:"name"`
	AvatarURL     string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type googleProvider struct {
	*oauth2.Config
}

func NewGoogleProvider(cfg config.GoogleAuth, redirectURL string) (AuthProvider, error) {
	return &googleProvider{
		Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  GoogleOauthAuthEndpoint,
				TokenURL: GoogleOauthTokenEndpoint,
			},
		},
	}, nil
}

func (g googleProvider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return g.Exchange(context.Background(), code)
}

func (g googleProvider) Name() string {
	return "google"
}

func (g googleProvider) GetUser(token *oauth2.Token) (*User, error) {

	var user GoogleUser

	if err := makeRequest(token, g.Config, GoogleUserAPI, &user); err != nil {
		return nil, err
	}

	data := &User{
		AuthenticationMethod: "google",
		Emails: Emails{{
			Email:    user.Email,
			Verified: user.EmailVerified,
			Primary:  true,
		}},
		Metadata: &Claims{
			Issuer:        GoogleUrl,
			Subject:       user.ID,
			Name:          user.Name,
			AvatarUrl:     user.AvatarURL,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
		},
	}

	return data, nil
}
