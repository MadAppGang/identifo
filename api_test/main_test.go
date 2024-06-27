package api_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/mail/mock"
	"gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/baloo.v3/assert"
)

type Config struct {
	ServerURL            string `env:"SERVER" envDefault:"http://localhost:8081"`
	AppID                string `env:"APP_ID" envDefault:"59fd884d8f6b180001f5b4e2"`
	AppSecret            string `env:"APP_SECRET" envDefault:"app_secret"`
	AppID2               string `env:"APP_ID2" envDefault:"59fd884d8f6b180001f5b4e3"`
	AppSecret2           string `env:"APP_SECRET2" envDefault:"app_secret_2"`
	RunTestServer        bool   `env:"RUN_TEST_SERVER" envDefault:"true"`
	User1                string `env:"USER1" envDefault:"test@madappgang.com"`
	User1Pswd            string `env:"USER1_PSWD" envDefault:"Secret3"`
	User2                string `env:"USER2" envDefault:"new_user@madappgang.com"`
	User2Pswd            string `env:"USER2_PSWD" envDefault:"Secret321"`
	User3                string `env:"USER3" envDefault:"new_user3@madappgang.com"`
	User3Pswd            string `env:"USER3_PSWD" envDefault:"Secret321"`
	User4                string `env:"USER4" envDefault:"new_user4@madappgang.com"`
	User4Pswd            string `env:"USER4_PSWD" envDefault:"Secret321_4"`
	ManagementKeyID1     string `env:"MANAGEMENT_KEY1" envDefault:"63c6273a46504e3abdc00fc7"`
	ManagementKeySecret1 string `env:"MANAGEMENT_KEY_SECRET1" envDefault:"secret1"`
	ManagementKeyID2     string `env:"MANAGEMENT_KEY1" envDefault:"63c6273a46504e3abdc00fc6"`
	ManagementKeySecret2 string `env:"MANAGEMENT_KEY_SECRET1" envDefault:"secret2"`
}

var cfg = Config{}

// test stores the HTTP testing client preconfigured
var request *baloo.Client

var emailService *mock.EmailService

// ============================================================
// Some test helper function here to setup test environment
// ============================================================
func TestMain(m *testing.M) {
	// forceIntegrationTests()

	// try to load dotenv file. If failed - just ignore. Dotenv file is optional
	_ = godotenv.Load()

	env.Parse(&cfg)
	request = baloo.New(cfg.ServerURL)

	var httpServer *http.Server
	if cfg.RunTestServer == true {
		_, httpServer = runServer()
	}

	code := m.Run()

	if cfg.RunTestServer == true {
		stopServer(httpServer)
	}

	os.Exit(code)
}

// run identifo server and import test data
func runServer() (model.Server, *http.Server) {
	var settings model.FileStorageSettings

	// if we do regular isolated tests - use boldtb as a storage
	if len(os.Getenv("IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION")) == 0 {
		os.Remove("../data/db.db")
		settings, _ = model.ConfigStorageSettingsFromString("file://../test/artifacts/api/config.yaml")
	} else {
		// if we do integration tests with mongodb - run tests with mongodb
		settings, _ = model.ConfigStorageSettingsFromString("file://../test/artifacts/api/config-mongo.yaml")
	}

	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil {
		log.Fatalf("Unable to load config with error: %v", err)
	}

	srv, err := config.NewServer(configStorage, nil)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	emailService = srv.Services().Email.Transport().(*mock.EmailService)

	if err := config.ImportApps("../test/artifacts/api/apps.json", srv.Storages().App, true); err != nil {
		log.Fatalf("error importing apps to server: %v", err)
	}
	if err := config.ImportUsers("../test/artifacts/api/users.json", srv.Storages().User, true); err != nil {
		log.Fatalf("error importing users to server: %v", err)
	}
	if err := config.ImportManagement("../test/artifacts/api/management_keys.json", srv.Storages().ManagementKey, true); err != nil {
		log.Fatalf("error importing management keys to server: %v", err)
	}

	// updates CORS after apps import
	srv.UpdateCORS()

	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}

	go func() {
		log.Println("web api ListenAndServe")
		if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// maybe move the te
	return srv, httpSrv
}

// stop server and clear the data
func stopServer(server *http.Server) {
	server.Close()
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
type (
	validatorFunc     = func(map[string]interface{}) error
	validatorFuncText = func(string) error
)

func validateJSON(validator validatorFunc) assert.Func {
	return func(res *http.Response, req *http.Request) error {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// Re-fill body reader stream after reading it
		res.Body = io.NopCloser(bytes.NewBuffer(body))

		// parse the data
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		return validator(data)
	}
}

func validateBodyText(validator validatorFuncText) assert.Func {
	return func(res *http.Response, req *http.Request) error {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// Re-fill body reader stream after reading it
		res.Body = io.NopCloser(bytes.NewBuffer(body))

		return validator(string(body))
	}
}

func Signature(data, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))

	if _, err := mac.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("error creating signature for data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func forceIntegrationTests() {
	os.Setenv("IDENTIFO_TEST_INTEGRATION", "1")
	os.Setenv("IDENTIFO_TEST_AWS_ENDPOINT", "http://localhost:9000")
	os.Setenv("AWS_ACCESS_KEY_ID", "testing")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "testing_secret")
	os.Setenv("IDENTIFO_FORCE_S3_PATH_STYLE", "1")
	os.Setenv("IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION", "1")
	os.Setenv("IDENTIFO_STORAGE_MONGO_CONN", "mongodb://admin:password@localhost:27017/billing-local?authSource=admin")
	os.Setenv("IDENTIFO_REDIS_HOST", "127.0.0.1:6379")
}
