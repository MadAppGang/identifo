package runner_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/madappgang/identifo/config"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/test/runner"

	"gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var request = baloo.New("http://localhost:8081")

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
		Status(200).
		Done()
}

// test wrong app ID login
func TestLoginWithWrongAppID(t *testing.T) {
	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "wrong_app_ID").
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
