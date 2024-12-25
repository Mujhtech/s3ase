package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/mujhtech/s3ase/config"
	"golang.org/x/oauth2"
)

const (
	GithubUrl                = "https://github.com"
	GithubApiBase            = "https://api.github.com"
	GithubUserAPI            = GithubApiBase + "/user"
	GithubEmailAPI           = GithubApiBase + "/user/emails"
	GithubOauthAuthEndpoint  = GithubUrl + "/login/oauth/authorize"
	GithubOauthTokenEndpoint = GithubUrl + "/login/oauth/access_token"
)

type GithubUser struct {
	ID        int    `json:"id"`
	UserName  string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GithubUserEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type githubProvider struct {
	*oauth2.Config
}

func NewGithubProvider(cfg config.GithubAuth, redirectURL string) (AuthProvider, error) {
	return &githubProvider{
		Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  GithubOauthAuthEndpoint,
				TokenURL: GithubOauthTokenEndpoint,
			},
		},
	}, nil
}

func (g githubProvider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return g.Exchange(context.Background(), code)
}

func (g githubProvider) Name() string {
	return "github"
}

func (g githubProvider) GetUser(token *oauth2.Token) (*User, error) {

	// https://github.com/teamhanko/hanko/blob/main/backend/thirdparty/provider_github.go

	var user GithubUser

	if err := makeRequest(token, g.Config, GithubUserAPI, &user); err != nil {
		return nil, err
	}

	data := &User{
		AuthenticationMethod: "github",
		Metadata: &Claims{
			Issuer:    GithubUrl,
			Subject:   strconv.Itoa(user.ID),
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
			Username:  user.UserName,
		},
	}

	var emails []*GithubUserEmail

	if err := makeRequest(token, g.Config, GithubEmailAPI, &emails); err != nil {
		return nil, err
	}

	for _, e := range emails {
		if e.Email != "" {
			data.Emails = append(data.Emails, Email{Email: e.Email, Verified: e.Verified, Primary: e.Primary})
		}

		if e.Primary {
			data.Metadata.Email = e.Email
			data.Metadata.EmailVerified = e.Verified
		}
	}

	if len(data.Emails) <= 0 {
		return nil, errors.New("unable to find email with GitHub provider")
	}

	return data, nil
}
