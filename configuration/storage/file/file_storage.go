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

// ConfigurationStorage is a wrapper over server configuration file.
type ConfigurationStorage struct {
	ServerConfigPath string
	UpdateChan       chan interface{}
	updateChanClosed bool
}

// NewConfigurationStorage creates and returns new file configuration storage.
func NewConfigurationStorage(settings model.ConfigurationStorageSettings) (*ConfigurationStorage, error) {
	return &ConfigurationStorage{
		ServerConfigPath: settings.SettingsKey,
		UpdateChan:       make(chan interface{}, 1),
	}, nil
}

// Insert writes new value to server configuration file.
func (cs *ConfigurationStorage) Insert(key string, value interface{}) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get server configuration file: %s", err)
	}

	if err = cs.updateConfigFile(value, filepath.Join(dir, key)); err != nil {
		return fmt.Errorf("Cannot update server configuration file: %s", err)
	}

	// Indicate config update. To prevent writing to closed channel, make a check.
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
func (cs *ConfigurationStorage) LoadServerSettings(ss *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get server configuration file: %s", err)
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, ss.ConfigurationStorage.SettingsKey))
	if err != nil {
		return fmt.Errorf("Cannot read server configuration file: %s", err)
	}

	if err = yaml.Unmarshal(yamlFile, ss); err != nil {
		return fmt.Errorf("Cannot unmarshal server configuration file: %s", err)
	}
	return nil
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

func (cs *ConfigurationStorage) updateConfigFile(in interface{}, dir string) error {
	ss, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("Cannot marshall configuration: %s", err)
	}

	if err = ioutil.WriteFile(dir, ss, 0644); err != nil {
		return fmt.Errorf("Cannot write configuration file: %s", err)
	}
	return nil
}
