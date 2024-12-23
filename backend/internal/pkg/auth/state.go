package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/pkg/encrypt"
)

type State struct {
	Provider   string    `json:"provider"`
	RedirectTo string    `json:"redirect_to"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Nonce      string    `json:"nonce"`
}

func GenerateState(config *config.Config, provider string, redirectTo string, options ...func(*State)) ([]byte, error) {
	if provider == "" {
		return nil, errors.New("provider must be present")
	}

	nonce, err := encrypt.GenerateRandomStringURLSafe(32)
	if err != nil {
		return nil, fmt.Errorf("could not generate nonce: %w", err)
	}

	now := time.Now().UTC()
	state := State{
		Provider:   provider,
		RedirectTo: redirectTo,
		IssuedAt:   now,
		ExpiresAt:  now.Add(time.Minute * 5),
		Nonce:      nonce,
	}

	for _, option := range options {
		option(&state)
	}

	stateJson, err := json.Marshal(state)

	if err != nil {
		return nil, fmt.Errorf("could not marshal state: %w", err)
	}

	aesCfb, err := encrypt.NewAesCfb(config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate aesgcm: %w", err)
	}

	encryptedState, err := aesCfb.Encrypt(stateJson)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt state: %w", err)
	}

	return []byte(encryptedState), nil
}

func VerifyState(config *config.Config, state string, expectedState string) (*State, error) {
	decodedState, err := decodeState(config, state)
	if err != nil {
		return nil, fmt.Errorf("could not decode state: %w", err)
	}

	decodedExpectedState, err := decodeState(config, expectedState)
	if err != nil {
		return nil, fmt.Errorf("could not decode expectedState: %w", err)
	}

	if decodedState.Nonce != decodedExpectedState.Nonce {
		return nil, errors.New("could not verify state")
	}

	if time.Now().UTC().After(decodedState.ExpiresAt) {
		return nil, errors.New("state is expired")
	}

	return decodedState, nil
}

func decodeState(config *config.Config, state string) (*State, error) {
	aesCfb, err := encrypt.NewAesCfb(config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate aesgcm: %w", err)
	}

	decryptedState, err := aesCfb.Decrypt(state)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt state: %w", err)
	}

	var unmarshalledState State
	err = json.Unmarshal([]byte(decryptedState), &unmarshalledState)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal state: %w", err)
	}
	return &unmarshalledState, nil
}
