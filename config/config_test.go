package config_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfigurationWithEmptyFlags(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	stor, err := config.InitConfigurationStorageFromFlag("")

	require.NoError(t, err)

	// the settings are not loaded yet, should be empty
	assert.Nil(t, stor.LoadedSettings())

	// load settings and check the settings are default and no errors
	settings, errs := stor.LoadServerSettings(true)
	assert.True(t, reflect.DeepEqual(settings, model.DefaultServerSettings))
	assert.Empty(t, errs)
}

func TestInitConfigurationWithFile(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	os.Mkdir("../data", os.ModePerm)
	stor, err := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_default_storage.yaml")
	require.NoError(t, err)

	// the settings are not loaded yet, should be empty
	assert.Nil(t, stor.LoadedSettings())

	// load settings and check the settings are default and no errors
	s, errs := stor.LoadServerSettings(true)
	// assert.True(t, reflect.DeepEqual(settings, model.DefaultServerSettings))
	assert.Empty(t, errs)
	assert.Equal(t, s.Storage.DefaultStorage.Type, model.DBTypeBoltDB)
	assert.Equal(t, s.Storage.DefaultStorage.BoltDB.Path, "../data/db.db")

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
	require.Empty(t, server.Errors())
	require.True(t, reflect.DeepEqual(s, server.Settings()))
	os.RemoveAll("../data")
}

func TestInitConfigurationWithWrongFilePath(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	stor, err := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_default_storage_not_exists.yaml")
	require.NoError(t, err)
	require.NotNil(t, stor)

	s, errs := stor.LoadServerSettings(true)
	assert.NotEmpty(t, errs)
	assert.Empty(t, s.Storage.DefaultStorage.Type) // empty storages

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestInitConfigurationWithWrongSettingsData(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	stor, err := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_wrong_data.yaml")
	require.NoError(t, err)
	require.NotNil(t, stor)

	s, errs := stor.LoadServerSettings(true)
	assert.NotEmpty(t, errs)
	assert.True(t, sliceContainsError(errs, "unsupported database type magicDB"))
	assert.Equal(t, string(s.Storage.DefaultStorage.Type), "magicDB") // wrong data

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
}

func sliceContainsError(errs []error, text string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error(), text) {
			return true
		}
	}
	return false
}

func TestInitConfigurationWithDefaultReferenceDefault(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	stor, err := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_default_reference_default.yaml")
	require.NoError(t, err)
	require.NotNil(t, stor)

	s, errs := stor.LoadServerSettings(true)
	assert.NotEmpty(t, errs)
	assert.True(t, sliceContainsError(errs, "DefaultStorage settings could not be of type Default"))
	assert.Equal(t, model.DBTypeDefault, s.Storage.DefaultStorage.Type) // default reference default

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestInitConfigurationWithBrokenSettingsAPICall(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	stor, _ := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_default_reference_default.yaml")
	s, errs := stor.LoadServerSettings(true)
	assert.NotEmpty(t, errs)
	assert.True(t, sliceContainsError(errs, "DefaultStorage settings could not be of type Default"))
	assert.Equal(t, model.DBTypeDefault, s.Storage.DefaultStorage.Type) // default reference default

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)

	// PING
	// let's make HTTP api call
	// Ping should be OK, even if configuration is wrong, indicating that server is running as expected
	req, _ := http.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Pong")

	// login request
	// the response should be error with settings are wrong
	req, _ = http.NewRequest("POST", "/auth/login", nil)
	rr = httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "DefaultStorage settings could not be of type Default") // actual error is there as well
	assert.Contains(t, rr.Body.String(), "Private key file not found")                           // keys unable to load as well

	// CORS request
	// the response should be allow for any app
	// to be able to get the config error from any location
	req, _ = http.NewRequest("OPTIONS", "/auth/app_settings", nil)
	rr = httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "DefaultStorage settings could not be of type Default") // actual error is there as well
	assert.Contains(t, rr.Body.String(), "Private key file not found")                           // keys unable to load as well
}

func TestInitConfigurationWithWithGoodConfigAndFailedStorages(t *testing.T) {
	os.Unsetenv(model.IdentifoConfigPathEnvName)
	// storages will failed to create as folder does not exists
	stor, err := config.InitConfigurationStorageFromFlag("file://../test/artifacts/configs/settings_with_wrong_file_path.yaml")
	require.NoError(t, err)

	// the settings are not loaded yet, should be empty
	assert.Nil(t, stor.LoadedSettings())

	// load settings and check the settings are default and no errors
	s, errs := stor.LoadServerSettings(true)
	// assert.True(t, reflect.DeepEqual(settings, model.DefaultServerSettings))
	assert.Empty(t, errs)
	assert.Equal(t, s.Storage.DefaultStorage.Type, model.DBTypeBoltDB)
	assert.Equal(t, s.Storage.DefaultStorage.BoltDB.Path, "/I/am/wrong/folder/db.db")

	server, err := config.NewServer(stor, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
	require.NotEmpty(t, server.Errors())
	require.True(t, reflect.DeepEqual(s, server.Settings()))

	// PING
	// let's make HTTP api call
	// Ping should be OK, even if configuration is wrong, indicating that server is running as expected
	req, _ := http.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Pong")

	// login request
	// the response should be error with settings are wrong
	req, _ = http.NewRequest("POST", "/auth/login", nil)
	rr = httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error creating app storage") // actual error is there as well
	assert.Contains(t, rr.Body.String(), "no such file or directory")
}
