package api_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/model"
	"gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/baloo.v3/assert"
)

type Config struct {
	// ServerURL     string `env:"SERVER" envDefault:"http://localhost:8081"`
	AppID         string `env:"APP_ID" envDefault:"59fd884d8f6b180001f5b4e2"`
	AppSecret     string `env:"APP_SECRET" envDefault:"app_secret"`
	RunTestServer bool   `env:"RUN_TEST_SERVER" envDefault:"true"`
	User1         string `env:"USER1" envDefault:"test@madappgang.com"`
	User1Pswd     string `env:"USER1_PSWD" envDefault:"Secret3"`
	User2         string `env:"USER2" envDefault:"new_user@madappgang.com"`
	User2Pswd     string `env:"USER2_PSWD" envDefault:"Secret321"`
	User3         string `env:"USER3" envDefault:"new_user3@madappgang.com"`
	User3Pswd     string `env:"USER3_PSWD" envDefault:"Secret321"`
	User4         string `env:"USER4" envDefault:"new_user4@madappgang.com"`
	User4Pswd     string `env:"USER4_PSWD" envDefault:"Secret321_4"`
}

var cfg = Config{}

// test stores the HTTP testing client preconfigured
var request *baloo.Client

// ============================================================
// Some test helper function here to setup test environment
// ============================================================
func TestMain(m *testing.M) {
	// try to load dotenv file. If failed - just ignore. Dotenv file is optional
	_ = godotenv.Load()

	env.Parse(&cfg)

	_, s := runServer()
	request = baloo.New(s.URL)

	code := m.Run()

	stopServer(s)

	os.Exit(code)
}

// run identifo server and import test data
func runServer() (model.Server, *httptest.Server) {
	var settings model.FileStorageSettings

	os.Remove("./db.db")
	settings, _ = model.ConfigStorageSettingsFromString("file://./config.yaml")

	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil {
		log.Fatalf("Unable to load config with error: %v", err)
	}

	srv, err := config.NewServer(configStorage, nil)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	log.Println("misconfig ListenAndServe")
	httpSrv := httptest.NewServer(srv.Router())

	return srv, httpSrv
}

// stop server and clear the data
func stopServer(server *httptest.Server) {
	server.Close()
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

func Signature(data, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))

	if _, err := mac.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("error creating signature for data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
