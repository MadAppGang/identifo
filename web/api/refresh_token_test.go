package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokens(t *testing.T) {
	ctx := testContext(testApp)

	reqBody := `{"scopes":["offline", "chat", "super_admin"]}`

	user := model.User{
		ID:         "test_user",
		Scopes:     []string{"chat"},
		Active:     true,
		Email:      "rt_some@example.com",
		Username:   "rt_some",
		Phone:      "1234567890",
		AccessRole: "user",
	}

	user, err := testServer.Storages().User.AddUserWithPassword(user, "qwerty", "user", false)
	require.NoError(t, err)

	tokenService := testServer.Services().Token

	refreshToken, err := tokenService.NewRefreshToken(
		user,
		[]string{"offline", "chat", "super_admin"},
		testApp)
	require.NoError(t, err)

	rts, err := tokenService.String(refreshToken)
	require.NoError(t, err)

	refreshToken, err = tokenService.Parse(rts)
	require.NoError(t, err)

	ctx = context.WithValue(ctx, model.TokenContextKey, refreshToken)
	ctx = context.WithValue(ctx, model.TokenRawContextKey, []byte(rts))

	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(reqBody))
	req = req.WithContext(ctx)

	rw := httptest.NewRecorder()

	h := testRouter.RefreshTokens()
	h(rw, req)

	require.Equal(t, http.StatusOK, rw.Code, rw.Body.String())

	c := claimsFromResponse(t, rw.Body.Bytes())
	assert.Equal(t, user.ID, c["sub"])
	assert.Equal(t, "test_app", c["aud"])
	assert.Equal(t, "chat offline", c["scopes"])

	c = refreshClaimsFromResponse(t, rw.Body.Bytes())
	assert.Equal(t, user.ID, c["sub"])
	assert.Equal(t, "test_app", c["aud"])
	assert.Equal(t, "chat offline", c["scopes"])
}
