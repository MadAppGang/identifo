package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/madappgang/identifo/v2/model"
	"gopkg.in/yaml.v2"
)

// ConfigurationStorage is a wrapper over server configuration file.
type ConfigurationStorage struct {
	ServerConfigPath string
	UpdateChan       chan interface{}
	updateChanClosed bool
	cache            *model.ServerSettings
	config           model.FileStorageSettings
	errors           []error
}

func NewDefaultConfigurationStorage() model.ConfigurationStorage {
	configPaths := []string{
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
				return c
			}
		}
	}
	// ok there is not default config files, the last line is to set up in right here in the code:
	return NewBuildingConfigurationStorage()
}

// NewConfigurationStorage creates and returns new file configuration storage.
func NewConfigurationStorage(config model.FileStorageSettings) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from specified file...")

	if config.Type != model.FileStorageTypeLocal {
		return nil, fmt.Errorf("could not create file config storage from non-file settings")
	}

	cs := &ConfigurationStorage{
		config:           config,
		ServerConfigPath: config.Local.Path,
		UpdateChan:       make(chan interface{}, 1),
	}

	return cs, nil
}

// WriteConfig writes new config to server configuration file.
func (cs *ConfigurationStorage) WriteConfig(settings model.ServerSettings) error {
	ss, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("Cannot marshall configuration: %s", err)
	}

	if err = os.WriteFile(cs.ServerConfigPath, ss, 0o644); err != nil {
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
func (cs *ConfigurationStorage) LoadServerSettings(validate bool) (model.ServerSettings, []error) {
	cs.errors = nil

	dir, err := os.Getwd()
	if err != nil {
		cs.errors = append(cs.errors, fmt.Errorf("Cannot get server configuration file: %s", err))
		return model.ServerSettings{}, cs.errors
	}

	yamlFile, err := os.ReadFile(filepath.Join(dir, cs.ServerConfigPath))
	if err != nil {
		cs.errors = append(cs.errors, fmt.Errorf("Cannot read server configuration file: %s", err))
		return model.ServerSettings{}, cs.errors
	}

	var settings model.ServerSettings
	if err = yaml.Unmarshal(yamlFile, &settings); err != nil {
		cs.errors = append(cs.errors, fmt.Errorf("Cannot unmarshal server configuration file: %s", err))
		return model.ServerSettings{}, cs.errors
	}

	settings.Config = cs.config
	cs.cache = &settings

	if validate {
		cs.errors = settings.Validate(true)
	}
	return settings, cs.errors
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

func (cs *ConfigurationStorage) LoadedSettings() *model.ServerSettings {
	return cs.cache
}

func (cs *ConfigurationStorage) Errors() []error {
	return cs.errors
}

// fileExists check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
