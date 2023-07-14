package controller

import (
	"context"
	"testing"

	"github.com/madappgang/identifo/v2/jwt/service"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvitationFromClaim(t *testing.T) {
	claims := map[string]any{
		"some claims":         "some value",
		"role:tenant1:group1": []any{model.RoleAdmin, model.RoleGuest},
		"role:tenant1:group2": []any{model.RoleOwner, model.RoleGuest},
		"role:tenant2:group1": []any{model.RoleGuest},
		"role:tenant2:group2": "some text",
	}

	inv := getInvitationFromClaim(claims)

	require.NotNil(t, inv)

	assert.NotNil(t, inv["tenant1"])
	assert.NotEmpty(t, inv["tenant1"].Groups)
	assert.Equal(t, 2, len(inv["tenant1"].Groups))
	assert.Len(t, inv["tenant1"].Groups["group1"], 2)
	assert.Len(t, inv["tenant1"].Groups["group2"], 2)
	assert.Contains(t, inv["tenant1"].Groups["group1"], model.RoleAdmin)
	assert.Contains(t, inv["tenant1"].Groups["group1"], model.RoleGuest)
	assert.Contains(t, inv["tenant1"].Groups["group2"], model.RoleOwner)

	assert.NotNil(t, inv["tenant2"])
	assert.NotEmpty(t, inv["tenant2"].Groups)
	assert.Equal(t, 1, len(inv["tenant2"].Groups))
	assert.Len(t, inv["tenant1"].Groups["group1"], 2)
}

func TestFilterInviteeCouldInvite(t *testing.T) {
	invitedTo := getInvitationFromClaim(claims)

	c := NewUserStorageController(nil, nil, nil, nil, nil, nil, nil, nil, model.DefaultServerSettings)
	r := c.filterInviteeCouldInvite(inviterMembership, invitedTo)
	require.NotNil(t, r)
	require.NotEmpty(t, r)

	assert.Len(t, r, 1)
	assert.Len(t, r["tenant1"].Groups, 1)
	assert.Len(t, r["tenant1"].Groups["group1"], 2)
	assert.Contains(t, r["tenant1"].Groups["group1"], model.RoleAdmin)
	assert.Contains(t, r["tenant1"].Groups["group1"], model.RoleGuest)
}

func TestAddUserWithInvitationToken(t *testing.T) {
	u := mock.UserStorage{
		UData: map[string]model.UserData{
			"user1": {
				UserID:           "user1",
				TenantMembership: inviterMembership,
			},
		},
	}
	user := model.User{
		ID:        "user1",
		GivenName: "Mr Inviter",
		Email:     "aooth@madappgang.com",
	}
	newUser := model.User{
		ID:        "New User",
		GivenName: "Mr New USer",
		Email:     "aooth_new@madappgang.com",
	}
	ts := createTokenService(t)
	token, err := ts.NewToken(model.TokenTypeInvite, user, nil, nil, claims)
	require.NoError(t, err)
	c := NewUserStorageController(&u, &u, nil, nil, nil, nil, nil, nil, model.DefaultServerSettings)

	c.AddUserToTenantWithInvitationToken(context.TODO(), newUser, token)
}

var claims = map[string]any{
	"some claims":         "some value",
	"role:tenant1:group1": []any{model.RoleAdmin, model.RoleGuest},
	"role:tenant1:group2": []any{model.RoleOwner, model.RoleGuest},
	"role:tenant2:group1": []any{model.RoleGuest},
	"role:tenant2:group2": "some text",
}

var inviterMembership = map[string]model.TenantMembership{
	"tenant1": {
		TenantID: "tenant1",
		Groups: map[string][]string{
			"group1": {model.RoleGuest, model.RoleAdmin},
			"group2": {model.RoleGuest},
		},
	},
	"tenant3": {
		TenantID: "tenant3",
		Groups:   map[string][]string{"group1": {model.RoleGuest, model.RoleAdmin}},
	},
}

const (
	keyPath    = "../../jwt/test_artifacts/private.pem"
	testIssuer = "aooth.madappgang.com"
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
