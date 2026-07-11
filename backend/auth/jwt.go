package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gotidy/ptr"
	"github.com/mujhtech/s3ase/cache"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

const (
	HeaderAuthorization = "Authorization"
	issuer              = "s3ase"

	userSessionTokenLifeTime time.Duration = 30 * 24 * time.Hour // 30 days.
)

type JWTAuth struct {
	cfg       *config.Config
	userRepo  store.UserRepository
	tokenRepo store.TokenRepository
	cache     cache.Cache
}

type Claims struct {
	gojwt.RegisteredClaims

	Value string `json:"value,omitempty"`

	Token *SubClaimsToken `json:"tkn,omitempty"`
}

type SubClaimsToken struct {
	Type CredentialType `json:"typ,omitempty"`
	ID   string         `json:"id,omitempty"`
}

func NewJWTAuth(cfg *config.Config, userRepo store.UserRepository, tokenRepo store.TokenRepository, cache cache.Cache) Auth {
	return &JWTAuth{
		cfg:       cfg,
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cache:     cache,
	}
}

func (j JWTAuth) UserAuthenticate(r *http.Request) (*UserSession, error) {

	ctx := r.Context()

	creds, err := getCredentials(r)

	if err != nil {
		return nil, err
	}

	if len(creds.Token) == 0 {
		return nil, errors.New("missing credentials")
	}

	claims, err := parseTokenClaims(creds.Token, j.cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("parsing of JWT claims failed: %w", err)
	}

	if claims.Value == "" {
		return nil, errors.New("jwt is missing subject value")
	}

	user, err := j.userRepo.FindUserByID(ctx, claims.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var metadata *TokenMetadata
	switch {
	case claims.Token != nil:
		metadata, err = j.metadataFromTokenClaims(ctx, user.ID, claims.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata from token claims: %w", err)
		}
	default:
		return nil, fmt.Errorf("jwt is missing sub-claims")
	}

	return &UserSession{
		User:     user,
		Metadata: metadata,
	}, nil

}

func parseTokenClaims(rawToken, secret string) (*Claims, error) {
	claims := &Claims{}
	parsed, err := gojwt.ParseWithClaims(rawToken, claims, func(_ *gojwt.Token) (interface{}, error) {
		return []byte(secret), nil
	},
		gojwt.WithValidMethods([]string{gojwt.SigningMethodHS256.Alg()}),
		gojwt.WithIssuer(issuer),
		gojwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("parsed JWT token is invalid")
	}
	return claims, nil
}

func (j JWTAuth) CreateToken(
	ctx context.Context,
	value string,
	tokenType string,
) (*models.Token, string, error) {
	issuedAt := time.Now()

	lifetime := ptr.Duration(userSessionTokenLifeTime)

	var expiresAt *int64
	if lifetime != nil {
		expiresAt = ptr.Int64(issuedAt.Add(*lifetime).UnixMilli())
	}

	// create db entry first so we get the id.
	token := models.Token{
		ID:        uuid.NewString(),
		Type:      tokenType,
		Value:     value,
		IssuedAt:  issuedAt.UnixMilli(),
		ExpiredAt: expiresAt,
	}

	if err := j.tokenRepo.CreateToken(ctx, &token); err != nil {
		return nil, "", fmt.Errorf("failed to store token in db: %w", err)
	}

	// create jwt token.
	jwtToken, err := generateToken(&token, j.cfg.EncryptionKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create jwt token: %w", err)
	}

	// set token in cache.
	if err := j.cache.Set(ctx, token.Value, token.ID, *lifetime); err != nil {
		return nil, "", fmt.Errorf("failed to set token in cache: %w", err)
	}

	return &token, jwtToken, nil
}

func getCredentials(r *http.Request) (*Credentials, error) {
	headerToken := r.Header.Get(HeaderAuthorization)

	switch {
	case strings.HasPrefix(headerToken, "Basic "):

		usr, pwd, _ := r.BasicAuth()

		return &Credentials{
			Type:     CredentialTypeBasic,
			Username: usr,
			Password: pwd,
		}, nil

	case strings.HasPrefix(headerToken, "Bearer "):

		return &Credentials{
			Type:  CredentialTypeBearer,
			Token: headerToken[7:],
		}, nil
	default:

		return nil, fmt.Errorf("unsupported authorization type")
	}
}

func generateToken(token *models.Token, secret string) (string, error) {
	var expiresAt *gojwt.NumericDate
	if token.ExpiredAt != nil {
		expiresAt = gojwt.NewNumericDate(time.UnixMilli(*token.ExpiredAt))
	}

	jwtToken := gojwt.NewWithClaims(gojwt.SigningMethodHS256, Claims{
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  gojwt.NewNumericDate(time.UnixMilli(token.IssuedAt)),
			ExpiresAt: expiresAt,
		},
		Value: token.Value,
		Token: &SubClaimsToken{
			Type: CredentialType(token.Type),
			ID:   token.ID,
		},
	})

	res, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return res, nil
}

func (j *JWTAuth) metadataFromTokenClaims(
	ctx context.Context,
	value string,
	tknClaims *SubClaimsToken,
) (*TokenMetadata, error) {
	// ensure tkn exists
	tkn, err := j.tokenRepo.FindTokenByID(ctx, tknClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find token in db: %w", err)
	}

	// protect against faked JWTs for other principals in case of single salt leak
	if value != tkn.Value {
		return nil, fmt.Errorf("JWT was for value %s while db token was for value %s",
			value, tkn.Value)
	}

	return &TokenMetadata{
		Type:     CredentialType(tkn.Type),
		Metadata: tkn.Metadata,
		TokenID:  tkn.ID,
	}, nil
}
