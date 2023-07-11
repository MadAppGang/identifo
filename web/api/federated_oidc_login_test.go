package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/api"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testApp = model.AppData{
	ID:     "test_app",
	Active: true,
	Type:   model.Web,
	OIDCSettings: model.OIDCSettings{
		ProviderName:    "test",
		ClientID:        "test",
		ClientSecret:    "test",
		EmailClaimField: "email",
	},
}

// test environment
var (
	testRouter *api.Router
	testServer model.Server
	oidcServer *httptest.Server
)

type testConfig struct {
	model.ConfigurationStorage
}

var testServerSettings = model.DefaultServerSettings

func (tc testConfig) LoadServerSettings(validate bool) (model.ServerSettings, []error) {
	testServerSettings.KeyStorage.Local.Path = "../../jwt/test_artifacts/private.pem"
	testServerSettings.Login.LoginWith.FederatedOIDC = true
	return testServerSettings, nil
}

func (tc testConfig) LoadedSettings() *model.ServerSettings {
	return &testServerSettings
}

func init() {
	var err error

	rc := make(chan bool, 1)
	testServer, err = config.NewServer(testConfig{}, rc)
	if err != nil {
		panic(err)
	}

	rs := api.RouterSettings{
		LoginWith: model.LoginWith{
			FederatedOIDC: true,
		},
		Server: testServer,
		Cors:   cors.New(model.DefaultCors),
	}

	testRouter, err = api.NewRouter(rs)
	if err != nil {
		panic(err)
	}

	oidcServer, _ = testOIDCServer()
}

func testContext(app model.AppData) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.AppDataContextKey, app)

	return ctx
}

func Test_Router_OIDCLogin_Redirect(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL
	ctx := testContext(testApp)

	redirect := "http://localhost:8080"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s", url.QueryEscape(redirect))

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler
	testRouter.OIDCLogin(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	expectedAuthURL := oidcServer.URL + "/auth?client_id=test&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=code&scope=openid&state=test"

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	locQuery := locURL.Query()
	assert.NotEmpty(t, locQuery.Get("state"))

	locQuery.Set("state", "test")
	locURL.RawQuery = locQuery.Encode()

	require.Equal(t, expectedAuthURL, locURL.String())

	// there should be cookie with session
	cookies := rw.Result().Cookies()
	assert.NotEmpty(t, cookies)
}

func Test_Router_OIDCLogin_Complete(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL
	ctx := testContext(testApp)

	redirect := "http://localhost:8080"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s", url.QueryEscape(redirect))

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler - should return redirect to auth
	testRouter.OIDCLogin(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	locQuery := locURL.Query()
	state := locQuery.Get("state")
	assert.NotEmpty(t, state)

	testCompleteURL := fmt.Sprintf("http://localhost:8081/auth/federated/oidc/complete?code=%s&state=%s", "test", state)
	r = httptest.NewRequest(http.MethodGet, testCompleteURL, nil)
	r = r.WithContext(ctx)
	r.Header.Set("Cookie", rw.Header().Get("Set-Cookie"))

	crw := httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(crw, r)

	require.Equal(t, http.StatusOK, crw.Code, "should return JWT token", crw.Body.String())

	c := claimsFromResponse(t, crw.Body.Bytes())

	assert.NotEmpty(t, c["sub"], c)
	assert.Equal(t, "test_app", c["aud"], c)

	crw = httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(crw, r)

	cc := claimsFromResponse(t, crw.Body.Bytes())
	assert.Equal(t, c["sub"], cc["sub"])
	assert.Equal(t, c["aud"], cc["aud"])
	assert.Equal(t, c["iss"], cc["iss"])
}

func Test_Router_OIDCLogin_Complete_ByEmail(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL
	ctx := testContext(testApp)

	us := testServer.Storages().User

	users, _, err := us.FetchUsers("", 0, 100)
	require.NoError(t, err)

	for _, v := range users {
		us.DeleteUser(v.ID)
	}

	newUser, err := us.AddUserWithPassword(model.User{
		ID:     "test_user",
		Email:  "some@example.com",
		Active: true,
	}, "qwerty", "admin", false)
	require.NoError(t, err)

	redirect := "http://localhost:8080"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s", url.QueryEscape(redirect))

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler - should return redirect to auth
	testRouter.OIDCLogin(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	locQuery := locURL.Query()
	state := locQuery.Get("state")
	assert.NotEmpty(t, state)

	testCompleteURL := fmt.Sprintf("http://localhost:8081/auth/federated/oidc/complete?code=%s&state=%s", "test", state)
	r = httptest.NewRequest(http.MethodGet, testCompleteURL, nil)
	r = r.WithContext(ctx)
	r.Header.Set("Cookie", rw.Header().Get("Set-Cookie"))

	crw := httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(crw, r)

	require.Equal(t, http.StatusOK, crw.Code, "should return JWT token", crw.Body.String())

	c := claimsFromResponse(t, crw.Body.Bytes())

	assert.Equal(t, newUser.ID, c["sub"], c)
	assert.Equal(t, "test_app", c["aud"], c)
}

func claimsFromResponse(t *testing.T, response []byte) jwt.MapClaims {
	var token map[string]any

	err := json.Unmarshal(response, &token)
	require.NoError(t, err)

	at := token["access_token"].(string)
	require.NotEmpty(t, at)

	c := jwt.MapClaims{}

	jp := jwt.Parser{}

	_, _, err = jp.ParseUnverified(at, c)
	require.NoError(t, err)

	return c
}

func Test_Router_OIDCLoginComplete_Routing(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL

	app := testServer.Storages().App
	a, err := app.CreateApp(testApp)
	require.NoError(t, err)

	// start oidc login
	redirect := "http://localhost:8080"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s&appId=%s",
		url.QueryEscape(redirect),
		a.ID)
	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	rw := httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusFound, rw.Result().StatusCode, rw.Body.String())

	// complete with app id in path
	testURL = "/auth/federated/oidc/complete/" + a.ID
	r = httptest.NewRequest(http.MethodGet, testURL, nil)
	rw = httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusBadRequest, rw.Result().StatusCode, rw.Body.String())
	require.Equal(t, "error.federated.code.error", errCode(t, rw.Body.Bytes()), rw.Body.String())

	// complete with app id in query
	testURL = "/auth/federated/oidc/complete?appId=" + a.ID
	r = httptest.NewRequest(http.MethodGet, testURL, nil)
	rw = httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusBadRequest, rw.Result().StatusCode, rw.Body.String())
	require.Equal(t, "error.federated.code.error", errCode(t, rw.Body.Bytes()), rw.Body.String())
}

func errCode(t *testing.T, respBody []byte) string {
	var resp map[string]any

	err := json.Unmarshal(respBody, &resp)
	require.NoError(t, err)

	errRespI, ok := resp["error"]
	if !ok {
		return ""
	}

	errResp, ok := errRespI.(map[string]any)
	if !ok {
		return ""
	}

	code, ok := errResp["id"].(string)
	if !ok {
		return ""
	}

	return code
}
