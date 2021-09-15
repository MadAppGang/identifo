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

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/madappgang/identifo/config"
	"github.com/madappgang/identifo/model"
	"gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/baloo.v3/assert"
)

type Config struct {
	ServerURL     string `env:"SERVER" envDefault:"http://localhost:8081"`
	AppID         string `env:"APP_ID" envDefault:"59fd884d8f6b180001f5b4e2"`
	AppSecret     string `env:"APP_SECRET" envDefault:"app_secret"`
	RunTestServer bool   `env:"RUN_TEST_SERVER" envDefault:"true"`
	User1         string `env:"USER1" envDefault:"test@madappgang.com"`
	User1Pswd     string `env:"USER1_PSWD" envDefault:"Secret3"`
	User2         string `env:"USER2" envDefault:"new_user@madappgang.com"`
	User2Pswd     string `env:"USER2_PSWD" envDefault:"Secret321"`
	User3         string `env:"USER3" envDefault:"new_user3@madappgang.com"`
	User3Pswd     string `env:"USER3_PSWD" envDefault:"Secret321"`
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
	request = baloo.New(cfg.ServerURL)

	var httpServer http.Server
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
func runServer() (model.Server, http.Server) {
	os.Remove("./db.db")
	settings, _ := model.ConfigStorageSettingsFromString("file://config.yaml")
	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil {
		log.Fatalf("Unable to load config with error: %v", err)
	}

	srv, err := config.NewServer(configStorage, nil)
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
