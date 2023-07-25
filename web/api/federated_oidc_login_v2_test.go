package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var testApp = model.AppData{
// 	ID:     "test_app",
// 	Active: true,
// 	Type:   model.Web,
// 	OIDCSettings: model.OIDCSettings{
// 		ProviderName:    "test",
// 		ClientID:        "test",
// 		ClientSecret:    "test",
// 		EmailClaimField: "email",
// 	},
// }

// // test environment
// var (
// 	testRouter *api.Router
// 	testServer model.Server
// 	oidcServer *httptest.Server
// )

// type testConfig struct {
// 	model.ConfigurationStorage
// }

// var testServerSettings = model.DefaultServerSettings

// func (tc testConfig) LoadServerSettings(validate bool) (model.ServerSettings, []error) {
// 	testServerSettings.KeyStorage.Local.Path = "../../jwt/test_artifacts/private.pem"
// 	testServerSettings.Login.LoginWith.FederatedOIDC = true
// 	return testServerSettings, nil
// }

// func (tc testConfig) LoadedSettings() *model.ServerSettings {
// 	return &testServerSettings
// }

// func init() {
// 	var err error

// 	rc := make(chan bool, 1)
// 	testServer, err = config.NewServer(testConfig{}, rc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rs := api.RouterSettings{
// 		LoginWith: model.LoginWith{
// 			FederatedOIDC: true,
// 		},
// 		Server: testServer,
// 		Cors:   cors.New(model.DefaultCors),
// 	}

// 	testRouter, err = api.NewRouter(rs)
// 	if err != nil {
// 		panic(err)
// 	}

// 	oidcServer, _ = testOIDCServer()
// }

// func testContext(app model.AppData) context.Context {
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, model.AppDataContextKey, app)

// 	return ctx
// }

func Test_Router_OIDCLoginV2_Redirect(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL
	ctx := testContext(testApp)

	redirect := "http://localhost:8080"
	state := "some_test_state"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s&state=%s", url.QueryEscape(redirect), state)

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler
	testRouter.OIDCLogin(true)(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	expectedAuthURL := oidcServer.URL + "/auth?client_id=test&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=code&scope=openid&state=some_test_state"

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	require.Equal(t, expectedAuthURL, locURL.String())

	// there should be no cookie with session
	cookies := rw.Result().Cookies()
	assert.Empty(t, cookies)
}

func Test_Router_OIDCLoginV2_Complete(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL
	ctx := testContext(testApp)

	redirect := "http://localhost:8080"
	state := "some_test_state"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s&state=%s", url.QueryEscape(redirect), state)

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler - should return redirect to auth
	testRouter.OIDCLogin(true)(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	locQuery := locURL.Query()
	assert.Equal(t, state, locQuery.Get("state"))

	testCompleteURL := fmt.Sprintf("http://localhost:8081/auth/federated/oidc/complete?code=%s", "test")
	r = httptest.NewRequest(http.MethodGet, testCompleteURL, nil)
	r = r.WithContext(ctx)

	crw := httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(false)(crw, r)

	require.Equal(t, http.StatusOK, crw.Code, "should return JWT token", crw.Body.String())

	c := claimsFromResponse(t, crw.Body.Bytes())

	assert.NotEmpty(t, c["sub"], c)
	assert.Equal(t, "test_app", c["aud"], c)

	crw = httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(false)(crw, r)

	cc := claimsFromResponse(t, crw.Body.Bytes())
	assert.Equal(t, c["sub"], cc["sub"])
	assert.Equal(t, c["aud"], cc["aud"])
	assert.Equal(t, c["iss"], cc["iss"])
}

func Test_Router_OIDCLogin_CompleteV2_ByEmail(t *testing.T) {
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
	state := "some_test_state"
	testURL := fmt.Sprintf("/auth/federated/oidc/login?redirectUrl=%s&state=%s", url.QueryEscape(redirect), state)

	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	r = r.WithContext(ctx)

	rw := httptest.NewRecorder()

	// call the test handler - should return redirect to auth
	testRouter.OIDCLogin(true)(rw, r)

	require.Equal(t, http.StatusFound, rw.Code, "should redirect to the provider", rw.Body.String())

	locURL, err := url.Parse(rw.Header().Get("Location"))
	require.NoError(t, err)

	locQuery := locURL.Query()
	assert.Equal(t, state, locQuery.Get("state"))

	testCompleteURL := fmt.Sprintf("http://localhost:8081/auth/federated/oidc/complete?code=%s", "test")
	r = httptest.NewRequest(http.MethodGet, testCompleteURL, nil)
	r = r.WithContext(ctx)

	crw := httptest.NewRecorder()
	// call complete handler - should return Identifo's JWT token
	testRouter.OIDCLoginComplete(false)(crw, r)

	require.Equal(t, http.StatusOK, crw.Code, "should return JWT token", crw.Body.String())

	c := claimsFromResponse(t, crw.Body.Bytes())

	assert.Equal(t, newUser.ID, c["sub"], c)
	assert.Equal(t, "test_app", c["aud"], c)
}

func Test_Router_OIDCLoginCompleteV2_Routing(t *testing.T) {
	testApp.OIDCSettings.ProviderURL = oidcServer.URL

	app := testServer.Storages().App
	a, err := app.CreateApp(testApp)
	require.NoError(t, err)

	// start oidc login
	redirect := "http://localhost:8080"
	testURL := fmt.Sprintf("/v2/auth/federated/oidc/login?redirectUrl=%s&appId=%s&state=test",
		url.QueryEscape(redirect),
		a.ID)
	r := httptest.NewRequest(http.MethodGet, testURL, nil)
	rw := httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusFound, rw.Result().StatusCode, rw.Body.String())

	// complete with app id in path
	testURL = "/v2/auth/federated/oidc/complete/" + a.ID
	r = httptest.NewRequest(http.MethodGet, testURL, nil)
	rw = httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusBadRequest, rw.Result().StatusCode, rw.Body.String())
	require.Equal(t, "error.federated.code.error", errCode(t, rw.Body.Bytes()), rw.Body.String())

	// complete with app id in query
	testURL = "/v2/auth/federated/oidc/complete?appId=" + a.ID
	r = httptest.NewRequest(http.MethodGet, testURL, nil)
	rw = httptest.NewRecorder()

	testRouter.ServeHTTP(rw, r)

	require.Equal(t, http.StatusBadRequest, rw.Result().StatusCode, rw.Body.String())
	require.Equal(t, "error.federated.code.error", errCode(t, rw.Body.Bytes()), rw.Body.String())
}
