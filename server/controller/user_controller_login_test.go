package controller_test

import (
	"context"
	"testing"

	"github.com/madappgang/identifo/v2/jwt/service"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server/controller"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTenantData(t *testing.T) {
	ud := model.UserData{
		TenantMembership: tm,
	}

	flattenData := controller.TenantData(ud.TenantMembership, nil)
	assert.Len(t, flattenData, 0)

	flattenData = controller.TenantData(ud.TenantMembership, []string{"tenant:all"})
	assert.Contains(t, flattenData, "role:tenant1:default")
	assert.Contains(t, flattenData, "role:tenant1:group1")
	assert.Contains(t, flattenData, "role:tenant2:default")
	assert.Contains(t, flattenData, "role:tenant2:group33")
	assert.Contains(t, flattenData, "tenant:tenant1")
	assert.Contains(t, flattenData, "tenant:tenant2")
	assert.Equal(t, flattenData["role:tenant1:default"], []string{"admin"})
	assert.Equal(t, flattenData["role:tenant1:group1"], []string{"user"})
	assert.Equal(t, flattenData["role:tenant2:default"], []string{"guest"})
	assert.Equal(t, flattenData["role:tenant2:group33"], []string{"admin"})
	assert.Equal(t, flattenData["tenant:tenant1"], "I am a tenant 1")
	assert.Equal(t, flattenData["tenant:tenant2"], "Apple corporation")
	assert.Len(t, flattenData, 6)

	scopes := []string{"tenant:tenant1", "tenant:tenant3"}
	flattenData = controller.TenantData(ud.TenantMembership, scopes)
	assert.Contains(t, flattenData, "role:tenant1:default")
	assert.Contains(t, flattenData, "role:tenant1:group1")
	assert.NotContains(t, flattenData, "role:tenant2:default")
	assert.NotContains(t, flattenData, "role:tenant2:group33")
	assert.Contains(t, flattenData, "tenant:tenant1")
	assert.NotContains(t, flattenData, "tenant:tenant2")
	assert.Equal(t, flattenData["role:tenant1:default"], []string{"admin"})
	assert.Equal(t, flattenData["role:tenant1:group1"], []string{"user"})
	assert.Equal(t, flattenData["tenant:tenant1"], "I am a tenant 1")
	assert.Len(t, flattenData, 3)

	scopes = []string{"tenant:tenant1", "tenant:tenant3", "tenant:all"}
	flattenData = controller.TenantData(ud.TenantMembership, scopes)
	assert.Len(t, flattenData, 6)
}

func TestAccessTokenScopes(t *testing.T) {
	scopes := []string{
		"id",
		"offline",
		"access:",
		"access:",
		"access:profile",
		"access:oidc",
	}
	r := controller.AccessTokenScopes(scopes)
	assert.Len(t, r, 2)
	assert.Contains(t, r, "profile")
	assert.Contains(t, r, "oidc")
}

func TestRequestJWT(t *testing.T) {
	u := mock.UserStorage{
		UData: map[string]model.UserData{
			"user1": {
				UserID:           "user1",
				TenantMembership: tm,
			},
		},
	}
	user := model.User{
		ID:    "user1",
		Email: "aooth@madappgang.com",
	}

	app := model.AppData{ID: "app1", Offline: true}
	tokenService := createTokenService(t)
	c := controller.NewUserStorageController(&u, nil, nil, nil, nil, nil, tokenService, nil, nil, model.ServerSettings{})

	scopes := []string{model.EmailScope}
	response, err := c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)

	scopes = append(scopes, model.OfflineScope)
	response, err = c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)

	access, err := tokenService.Parse(*response.AccessToken)
	require.NoError(t, err)

	// access token should have not email, only id token
	assert.NotContains(t, access.FullClaims().Payload, "email")
	assert.NotContains(t, access.FullClaims().Payload, "tenant:tenant1")

	// lets create a token
	assert.Nil(t, response.IDToken)
	scopes = append(scopes, model.IDTokenScope)
	response, err = c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	require.NotEmpty(t, response.IDToken)

	// id token should have email data but not tenant yet
	id, err := tokenService.Parse(*response.IDToken)
	require.NoError(t, err)
	assert.Contains(t, id.FullClaims().Payload, "email")
	assert.Equal(t, id.FullClaims().Payload["email"], user.Email)
	assert.NotContains(t, id.FullClaims().Payload, "tenant:tenant1")

	// let's add email to access token
	scopes = append(scopes, model.AccessTokenScopePrefix+model.EmailScope)
	response, err = c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	access, err = tokenService.Parse(*response.IDToken)
	require.NoError(t, err)
	assert.Contains(t, access.FullClaims().Payload, "email")
	assert.Equal(t, access.FullClaims().Payload["email"], user.Email)
	assert.NotContains(t, access.FullClaims().Payload, "tenant:tenant1")

	// add specific tenant membership scope to access token
	scopes = append(scopes, model.AccessTokenScopePrefix+model.TenantScopePrefix+"tenant1")
	response, err = c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	access, err = tokenService.Parse(*response.AccessToken)
	require.NoError(t, err)
	assert.Contains(t, access.FullClaims().Payload, "tenant:tenant1")
	assert.Equal(t, "I am a tenant 1", access.FullClaims().Payload["tenant:tenant1"])
	// we have information only for tenant1, nothing about tenant2
	assert.NotContains(t, access.FullClaims().Payload, "tenant:tenant2")

	// not let's get all tenatats in token
	scopes = []string{model.AccessTokenScopePrefix + model.TenantScopeAll}
	response, err = c.GetJWTTokens(context.TODO(), app, user, scopes)
	require.NoError(t, err)
	access, err = tokenService.Parse(*response.AccessToken)
	require.NoError(t, err)
	assert.Contains(t, access.FullClaims().Payload, "tenant:tenant1")
	assert.Equal(t, "I am a tenant 1", access.FullClaims().Payload["tenant:tenant1"])
	// we have information only for tenant1, nothing about tenant2
	assert.Contains(t, access.FullClaims().Payload, "tenant:tenant2")
	assert.Equal(t, tm["tenant2"].TenantName, access.FullClaims().Payload["tenant:tenant2"])
	assert.Equal(t, []any{"guest"}, access.FullClaims().Payload["role:tenant2:default"])
}

const (
	keyPath    = "../../jwt/test_artifacts/private.pem"
	testIssuer = "aooth.madappgang.com"
)

var tm = map[string]model.TenantMembership{
	"tenants1": {
		TenantID:   "tenant1",
		TenantName: "I am a tenant 1",
		Groups:     map[string][]string{"default": {"admin"}, "group1": {"user"}},
	},
	"tenant2": {
		TenantID:   "tenant2",
		TenantName: "Apple corporation",
		Groups:     map[string][]string{"default": {"guest"}, "group33": {"admin"}},
	},
}

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
