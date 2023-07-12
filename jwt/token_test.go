package jwt_test

import (
	"reflect"
	"testing"

	jwt "github.com/madappgang/identifo/v2/jwt/service"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keyPath            = "./test_artifacts/private.pem"
	testIssuer         = "identifo.madappgang.com"
	tokenStringExample = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsInN1YiI6IjEyMzQ1Njc4OTAifQ.AAlGn8m8YG3emPa8CIS6TS-ndqaZCGUydnhU8FznyZ1McYQKkLlcqDW2c04q9ZxKDZHeiSyNIDOKA-EP0GVthQ"
)

// TODO: refactor new storage type

func createTokenService(t *testing.T) model.TokenService {
	us, err := mem.NewUserStorage()
	if err != nil {
		t.Fatalf("Unable to create user storage %v", err)
	}
	tstor, err := mem.NewTokenStorage()
	if err != nil {
		t.Fatalf("Unable to create token storage %v", err)
	}
	as, err := mem.NewAppStorage()
	if err != nil {
		t.Fatalf("Unable to create app storage %v", err)
	}

	keyStorage, err := storage.NewKeyStorage(model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			Path: keyPath,
		},
	})
	require.NoError(t, err)

	privateKey, err := keyStorage.LoadPrivateKey()
	require.NoError(t, err)

	tokenService, err := jwt.NewJWTokenService(
		privateKey,
		testIssuer,
		tstor,
		as,
		us,
	)
	require.NoError(t, err)
	return tokenService
}

func TestParseString(t *testing.T) {
	tokenService := createTokenService(t)

	token, err := tokenService.Parse(tokenStringExample)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, ok := token.(*model.JWToken)
	assert.True(t, ok)

	assert.Equal(t, token.Subject(), "1234567890")
	assert.Equal(t, token.IssuedAt().Unix(), int64(1516239022))
}

func TestTokenToString(t *testing.T) {
	tokenService := createTokenService(t)

	token, err := tokenService.Parse(tokenStringExample)
	require.NoError(t, err)
	require.NotNil(t, token)

	tokenString, err := tokenService.String(token)
	assert.NoError(t, err)

	token2, err := tokenService.Parse(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, token2)

	t1, _ := token.(*model.JWToken)
	t2, _ := token2.(*model.JWToken)
	claims1 := t1.Claims
	claims2 := t1.Claims

	if !reflect.DeepEqual(t1.Header, t2.Header) {
		t.Errorf("Headers = %+v, want %+v", t1.Header, t2.Header)
	}
	if !reflect.DeepEqual(claims1, claims2) {
		t.Errorf("Claims = %+v, want %+v", claims1, claims2)
	}
}

func TestNewToken(t *testing.T) {
	tokenService := createTokenService(t)

	ustg, _ := mem.NewUserStorage()
	user, _ := ustg.AddUserWithPassword(model.User{
		ID:       "12345566",
		Username: "username",
		Email:    "username@gmailc.om",
	}, "password", "admin", false)
	scopes := []string{"scope1", "scope2"}
	app := model.AppData{
		ID:                           "123456",
		Secret:                       "1",
		Active:                       true,
		Name:                         "testName",
		Description:                  "testDescriprion",
		Scopes:                       scopes,
		Offline:                      true,
		Type:                         model.Web,
		RedirectURLs:                 []string{},
		RegistrationForbidden:        false,
		AnonymousRegistrationAllowed: true,
		NewUserDefaultRole:           "",
	}
	token, err := tokenService.NewAccessToken(user, scopes, app, false, nil)
	assert.NoError(t, err)

	tokenString, err := tokenService.String(token)
	assert.NoError(t, err)

	token2, err := tokenService.Parse(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, token2)

	t2, _ := token2.(*model.JWToken)
	assert.NotNil(t, t2.Payload()["name"])
	assert.Equal(t, testIssuer, t2.Issuer())
	assert.Equal(t, user.ID, t2.Subject())
	assert.Equal(t, app.ID, t2.Audience())
}
