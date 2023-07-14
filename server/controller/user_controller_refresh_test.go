package controller_test

import (
	"context"
	"testing"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server/controller"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshJWTToken(t *testing.T) {
	user := model.User{
		ID:    "user1",
		Email: "user@email.com",
	}

	u := mock.UserStorage{
		UData: map[string]model.UserData{
			"user1": {
				UserID:           "user1",
				TenantMembership: tm,
			},
		},
		Users: []model.User{user},
	}
	app := model.AppData{ID: "app1", Offline: true}
	tokenService := createTokenService(t)
	tokenStorage := mock.NewTokenStorage()
	c := controller.NewUserStorageController(&u, nil, nil, nil, tokenStorage, nil, tokenService, nil, nil, model.ServerSettings{})

	token, err := tokenService.NewToken(model.TokenTypeAccess, user, []string{"app1"}, []string{"Email"}, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, token.Raw)

	refresh, err := tokenService.NewToken(model.TokenTypeRefresh, user, []string{"app1"}, nil, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, token.Raw)

	response, err := c.RefreshJWTToken(context.TODO(), refresh, token.Raw, app, []string{model.EmailScope, model.IDTokenScope, model.OfflineScope})
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotEmpty(t, response.IDToken)

	// let's ensure we could not use the same refresh token again
	_, err = c.RefreshJWTToken(context.TODO(), refresh, token.Raw, app, []string{model.EmailScope, model.IDTokenScope, model.OfflineScope})
	require.Error(t, err)
	require.ErrorIs(t, err, l.ErrorTokenBlocked)
}
