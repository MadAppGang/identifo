package jwt_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/jwt/service"
	"github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keyPath            = "./test_artifacts/private.pem"
	testIssuer         = "aooth.madappgang.com"
	tokenStringExample = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsInN1YiI6IjEyMzQ1Njc4OTAifQ.AAlGn8m8YG3emPa8CIS6TS-ndqaZCGUydnhU8FznyZ1McYQKkLlcqDW2c04q9ZxKDZHeiSyNIDOKA-EP0GVthQ"
)

func createTokenService(t *testing.T) model.TokenService {
	keyStorage, err := storage.NewKeyStorage(model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			Path: keyPath,
		},
	})
	require.NoError(t, err)

	privateKey, err := keyStorage.LoadPrivateKey()
	require.NoError(t, err)

	tokenService, err := service.NewJWTokenService(
		privateKey,
		testIssuer,
		model.DefaultServerSettings.SecuritySettings,
	)
	require.NoError(t, err)
	return tokenService
}

func TestJWTSerialization(t *testing.T) {
	u := model.User{
		ID:    "user1",
		Email: "user@aooth.com",
	}
	data := map[string]any{
		"tenant:t1":       "MadAppGang",
		"role:t1:default": "owner",
		"role:t1:aooth":   "user",
	}
	ts := createTokenService(t)

	token, err := ts.NewToken(model.TokenTypeAccess, u, []string{"app1"}, []string{"Email"}, data)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := token.Claims
	clJson, err := json.Marshal(claims)
	require.NoError(t, err)
	require.NotEmpty(t, clJson)

	fmt.Println(string(clJson))

	tokenStr, err := ts.SignToken(token)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
	fmt.Println(tokenStr)

	tkn, err := ts.Parse(tokenStr)
	assert.NoError(t, err)
	assert.NotEmpty(t, tkn)

	// empty validator with default settings
	vld := validator.NewValidator(nil, nil, nil, nil)
	err = vld.Validate(token)
	assert.NoError(t, err)

	vld = validator.NewValidator([]string{"app1", "app2"}, nil, nil, nil)
	err = vld.Validate(token)
	assert.NoError(t, err)
}

func TestJWTValidation(t *testing.T) {
	u := model.User{
		ID:    "user1",
		Email: "user@aooth.com",
	}
	data := map[string]any{
		"tenant:t1":       "MadAppGang",
		"role:t1:default": "owner",
		"role:t1:aooth":   "user",
	}
	ts := createTokenService(t)

	token, err := ts.NewToken(model.TokenTypeAccess, u, []string{"app1"}, []string{"Email"}, data)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// most relaxed config
	config := validator.Config{
		Audience:           nil,
		Issuer:             nil,
		UserID:             nil,
		TokenType:          nil,
		PublicKey:          ts.PublicKey(),
		IsAudienceRequired: false,
		IsIssuerRequired:   false,
	}

	// empty validator with default settings
	vld, err := validator.NewValidatorWithConfig(config)
	require.NoError(t, err)
	err = vld.Validate(token)
	assert.NoError(t, err)

	vld = validator.NewValidator([]string{"app1", "app2"}, nil, nil, nil)
	err = vld.Validate(token)
	assert.NoError(t, err)

	vld = validator.NewValidator([]string{"app3", "app2"}, nil, nil, nil)
	err = vld.Validate(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error.validation.token.invalid.audience")

	config.TokenType = []string{string(model.TokenTypeAccess)}
	vld, err = validator.NewValidatorWithConfig(config)
	require.NoError(t, err)
	err = vld.Validate(token)
	assert.NoError(t, err)

	config.TokenType = []string{string(model.TokenTypeRefresh)}
	vld, err = validator.NewValidatorWithConfig(config)
	require.NoError(t, err)
	err = vld.Validate(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error.validation.token.invalid.type")
}
