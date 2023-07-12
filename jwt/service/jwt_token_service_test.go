package service_test

import (
	"reflect"
	"testing"

	"github.com/madappgang/identifo/v2/jwt/service"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keyPath            = "../test_artifacts/private.pem"
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

func TestParseString(t *testing.T) {
	tokenService := createTokenService(t)

	token, err := tokenService.Parse(tokenStringExample)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// claims
	_, ok := token.Claims.(*model.Claims)
	require.True(t, ok)

	// assert.Equal(t, string(model.TokenTypeAccess), token.Type())

	assert.Equal(t, "1234567890", token.Subject())
	assert.Equal(t, int64(1516239022), token.IssuedAt().Unix())
}

func TestTokenToString(t *testing.T) {
	tokenService := createTokenService(t)

	token, err := tokenService.Parse(tokenStringExample)
	require.NoError(t, err)
	require.NotNil(t, token)

	tokenString, err := tokenService.SignToken(token)
	assert.NoError(t, err)

	token2, err := tokenService.Parse(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, token2)

	claims1 := token.Claims
	claims2 := token2.Claims

	if !reflect.DeepEqual(token.Header, token2.Header) {
		t.Errorf("Headers = %+v, want %+v", token.Header, token2.Header)
	}
	if !reflect.DeepEqual(claims1, claims2) {
		t.Errorf("Claims = %+v, want %+v", claims1, claims2)
	}
}

func TestNewToken(t *testing.T) {
	tokenService := createTokenService(t)

	user := model.User{
		ID:       "12345566",
		Username: "username",
		Email:    "username@gmailc.om",
	}
	token, err := tokenService.NewToken(model.TokenTypeAccess, user, []string{"12345"}, []string{"Email"}, nil)
	assert.NoError(t, err)

	tokenString, err := tokenService.SignToken(token)
	assert.NoError(t, err)

	_, err = tokenService.Parse(tokenString)
	require.NoError(t, err)
}
