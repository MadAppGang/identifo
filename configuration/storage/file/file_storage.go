package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	keyStorageLocal "github.com/madappgang/identifo/configuration/key_storage/local"
	keyStorageS3 "github.com/madappgang/identifo/configuration/key_storage/s3"
	ijwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"gopkg.in/yaml.v2"
)

// config path for server name
const serverConfigPathEnvName = "SERVER_CONFIG_PATH"

// ConfigurationStorage is a wrapper over server configuration file.
type ConfigurationStorage struct {
	ServerConfigPath string
	keyStorage       model.KeyStorage
	UpdateChan       chan interface{}
	updateChanClosed bool
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
			c, e := NewConfigurationStorage(p)
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
func NewConfigurationStorage(config string) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from specified file...")
	filename := config
	if strings.HasPrefix(strings.ToUpper(filename), "FILE://") {
		filename = filename[7:]
	}

	cs := &ConfigurationStorage{
		ServerConfigPath: filename,
		UpdateChan:       make(chan interface{}, 1),
	}

	settings := model.ServerSettings{}
	if err := cs.LoadServerSettings(&settings); err != nil {
		return nil, fmt.Errorf("Cannot not load settings from local file config storage: %s", err)
	}

	var err error
	var keyStorage model.KeyStorage
	switch settings.ConfigurationStorage.KeyStorage.Type {
	case model.KeyStorageTypeLocal:
		keyStorage, err = keyStorageLocal.NewKeyStorage(settings.ConfigurationStorage.KeyStorage)
	case model.KeyStorageTypeS3:
		keyStorage, err = keyStorageS3.NewKeyStorage(settings.ConfigurationStorage.KeyStorage)
	default:
		return nil, fmt.Errorf("Unknown key storage type: %s", settings.ConfigurationStorage.KeyStorage.Type)
	}
	if err != nil {
		return nil, err
	}
	cs.keyStorage = keyStorage
	return cs, nil
}

// InsertConfig writes new value to server configuration file.
func (cs *ConfigurationStorage) InsertConfig(key string, value interface{}) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get server configuration file: %s", err)
	}

	if err = cs.updateConfigFile(value, filepath.Join(dir, key)); err != nil {
		return fmt.Errorf("Cannot update server configuration file: %s", err)
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
func (cs *ConfigurationStorage) LoadServerSettings(ss *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get server configuration file: %s", err)
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, cs.ServerConfigPath))
	if err != nil {
		return fmt.Errorf("Cannot read server configuration file: %s", err)
	}

	if err = yaml.Unmarshal(yamlFile, ss); err != nil {
		return fmt.Errorf("Cannot unmarshal server configuration file: %s", err)
	}
	return ss.Validate()
}

// InsertKeys inserts new public and private keys.
func (cs *ConfigurationStorage) InsertKeys(keys *model.JWTKeys) error {
	if err := cs.keyStorage.InsertKeys(keys); err != nil {
		return err
	}
	return nil
}

// LoadKeys loads public and private keys from the key storage.
func (cs *ConfigurationStorage) LoadKeys(alg ijwt.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
	return cs.keyStorage.LoadKeys(alg)
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

// fileExists check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
