package runner_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/madappgang/identifo/config"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/test/runner"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/baloo.v3/assert"
)

// test stores the HTTP testing client preconfigured
var request = baloo.New("http://localhost:8081")

// ============================================================
// Login with username and password
// ============================================================

// test happy day login
func TestLogin(t *testing.T) {
	data := `
	{
		"username": "test@madappgang.com",
		"password": "Secret3",
		"scopes": ["offline"]
	}`
	signature, _ := runner.Signature(data, "app_secret")

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "59fd884d8f6b180001f5b4e2").
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_with_refresh_scheme.json").
		Done()
}

// test happy day login, with no refresh token included
func TestLoginWithNoRefresh(t *testing.T) {
	g := NewGomegaWithT(t)

	data := `
	{
		"username": "test@madappgang.com",
		"password": "Secret3",
		"scopes": []
	}`
	signature, _ := runner.Signature(data, "app_secret")

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "59fd884d8f6b180001f5b4e2").
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data).To(MatchKeys(IgnoreExtras|IgnoreMissing, Keys{
				"access_token":  Not(BeZero()),
				"refresh_token": BeZero(),
			}))
			return nil
		})).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_scheme.json").
		Done()
}

// test wrong app ID login
func TestLoginWithWrongAppID(t *testing.T) {
	g := NewGomegaWithT(t)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "wrong_app_ID").
		Expect(t).
		AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data["error"]).To(MatchAllKeys(Keys{
				"id":               Equal("error.api.request.app_id.invalid"),
				"message":          Not(BeZero()),
				"detailed_message": Not(BeZero()),
				"status":           BeNumerically("==", 400),
			}))
			return nil
		})).
		Type("json").
		Status(400).
		Done()
}

// test wrong signature for mobile app
func TestLoginWithWrongSignature(t *testing.T) {
	data := `
	{
		"username": "test@madappgang.com",
		"password": "Secret3",
		"scopes": ["offline"]
	}`
	signature, _ := runner.Signature(data, "app_secret")

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "wrong_app_ID").
		SetHeader("Digest", "SHA-256="+signature+"_wrong").
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		Status(400).
		Done()
}

// ============================================================
// Some test helper function here to setup test environment
// ============================================================
func TestMain(m *testing.M) {
	_, httpServer := runServer()
	code := m.Run()
	stopServer(httpServer)
	os.Exit(code)
}

// run identifo server and import test data
func runServer() (model.Server, http.Server) {
	os.Remove("./db.db")
	settings, _ := model.ConfigStorageSettingsFromString("file://config.yaml")
	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil {
		log.Fatalf("Unable to load config with error: %v", err)
	}

	srv, err := config.NewServer(configStorage)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	if err := config.ImportApps("data/apps.json", srv.Storages().App); err != nil {
		log.Fatalf("error importing apps to server: %v", err)
	}
	if err := config.ImportUsers("data/users.json", srv.Storages().User); err != nil {
		log.Fatalf("error importing users to server: %v", err)
	}

	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// maybe move the te
	return srv, *httpSrv
}

// stop server and clear the data
func stopServer(server http.Server) {
	server.Shutdown(context.Background())
	log.Println("the server is gracefully stopped, bye 👋")
	log.Println("Stop server")
}

// Helper functions

// DumpResponse
func dumpResponse(res *http.Response, req *http.Request) error {
	data, err := httputil.DumpResponse(res, true)
	fmt.Printf("Response: %s \n", string(data))
	return err
}

// DumpResponse
type validatorFunc = func(map[string]interface{}) error

func validateJSON(validator validatorFunc) assert.Func {
	return func(res *http.Response, req *http.Request) error {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// Re-fill body reader stream after reading it
		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// parse the data
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		return validator(data)
	}
}
