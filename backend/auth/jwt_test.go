package auth

import (
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-only-32-byte-signing-secret!"

func TestGenerateAndParseTokenClaims(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	expiresAt := now.Add(time.Hour).UnixMilli()
	token := &models.Token{
		ID: "token-id", Value: "user-id", Type: string(CredentialTypeBearer),
		IssuedAt: now.UnixMilli(), ExpiredAt: &expiresAt,
	}

	rawToken, err := generateToken(token, testJWTSecret)
	require.NoError(t, err)

	claims, err := parseTokenClaims(rawToken, testJWTSecret)
	require.NoError(t, err)
	require.Equal(t, issuer, claims.Issuer)
	require.Equal(t, token.Value, claims.Value)
	require.Equal(t, token.ID, claims.Token.ID)
	require.Equal(t, CredentialTypeBearer, claims.Token.Type)
	require.Equal(t, now.Unix(), claims.IssuedAt.Unix())
	require.Equal(t, expiresAt/1000, claims.ExpiresAt.Unix())
}

func TestParseTokenClaimsRejectsUnexpectedAlgorithm(t *testing.T) {
	rawToken, err := gojwt.NewWithClaims(gojwt.SigningMethodHS384, Claims{
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer: issuer, ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Value: "user-id",
	}).SignedString([]byte(testJWTSecret))
	require.NoError(t, err)

	_, err = parseTokenClaims(rawToken, testJWTSecret)
	require.Error(t, err)
}
