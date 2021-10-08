package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/madappgang/identifo/model"
	"gopkg.in/yaml.v2"
)

// config path for server name
const serverConfigPathEnvName = "SERVER_CONFIG_PATH"

// ConfigurationStorage is a wrapper over server configuration file.
type ConfigurationStorage struct {
	ServerConfigPath string
	UpdateChan       chan interface{}
	updateChanClosed bool
	cache            model.ServerSettings
	cached           bool
	config           model.ConfigStorageSettings
}

func NewDefaultConfigurationStorage() (*ConfigurationStorage, error) {
	configPaths := []string{
		os.Getenv(serverConfigPathEnvName),
		"./server-config.yaml",
		"../../server/server-config.yaml",
	}

	for _, p := range configPaths {
		if p == "" {
			continue
		}
		if fileExists(p) {
			cs, _ := model.ConfigStorageSettingsFromStringFile(p)
			c, e := NewConfigurationStorage(cs)
			// if error, trying to other file from the list
			if e != nil {
				log.Printf("Unable to load default config from file %s, trying other one from the list (if any)", p)
				continue
			} else {
				log.Printf("Successfully loaded default config from  file %s", p)
				return c, nil
			}
		}
	}
	err := fmt.Errorf("Unable to load default config file from the following candidates: %+v", configPaths)
	log.Println(err)
	return nil, err
}

// NewConfigurationStorage creates and returns new file configuration storage.
func NewConfigurationStorage(config model.ConfigStorageSettings) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from specified file...")

	if config.Type != model.ConfigStorageTypeFile {
		return nil, fmt.Errorf("cold not create file config storage from non-file settings")
	}

	cs := &ConfigurationStorage{
		config:           config,
		ServerConfigPath: config.File.FileName,
		UpdateChan:       make(chan interface{}, 1),
	}
	log.Printf("Successfully loaded config data from %s\n", config.File.FileName)

	return cs, nil
}

// WriteConfig writes new config to server configuration file.
func (cs *ConfigurationStorage) WriteConfig(settings model.ServerSettings) error {
	ss, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("Cannot marshall configuration: %s", err)
	}

	if err = ioutil.WriteFile(cs.ServerConfigPath, ss, 0644); err != nil {
		return fmt.Errorf("Cannot write configuration file: %s", err)
	}

	// Indicate config update. To prevent writing to a closed channel, make a check.
	go func() {
		if cs.updateChanClosed {
			log.Println("Attempted to write to closed UpdateChan")
			return
		}
		cs.UpdateChan <- struct{}{}
	}()
	return nil
}

// LoadServerSettings loads server settings from the file.
func (cs *ConfigurationStorage) LoadServerSettings(forceReload bool) (model.ServerSettings, error) {
	if !forceReload && cs.cached {
		return cs.cache, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return model.ServerSettings{}, fmt.Errorf("Cannot get server configuration file: %s", err)
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, cs.ServerConfigPath))
	if err != nil {
		return model.ServerSettings{}, fmt.Errorf("Cannot read server configuration file: %s", err)
	}

	var settings model.ServerSettings
	if err = yaml.Unmarshal(yamlFile, &settings); err != nil {
		return model.ServerSettings{}, fmt.Errorf("Cannot unmarshal server configuration file: %s", err)
	}

	settings.Config = cs.config
	cs.cache = settings
	cs.cached = true

	return settings, settings.Validate()
}

// GetUpdateChan returns update channel.
func (cs *ConfigurationStorage) GetUpdateChan() chan interface{} {
	return cs.UpdateChan
}

// CloseUpdateChan closes update channel.
func (cs *ConfigurationStorage) CloseUpdateChan() {
	close(cs.UpdateChan)
	cs.updateChanClosed = true
}

func (cs *ConfigurationStorage) ForceReloadOnWriteConfig() bool {
	return false
}

// fileExists check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
