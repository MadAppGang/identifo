package s3

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	keyStorageLocal "github.com/madappgang/identifo/config/key_storage/local"
	keyStorageS3 "github.com/madappgang/identifo/config/key_storage/s3"
	"github.com/madappgang/identifo/model"
	s3Storage "github.com/madappgang/identifo/storage/s3"
	"gopkg.in/yaml.v2"
)

// ConfigurationStorage is a server configuration storage in S3.
type ConfigurationStorage struct {
	Client           *s3.S3
	Bucket           string
	ObjectName       string
	UpdateChan       chan interface{}
	updateChanClosed bool
	keyStorage       model.KeyStorage
	config           model.ConfigStorageSettings
	cache            model.ServerSettings
	cached           bool
}

// NewConfigurationStorage creates new server config storage in S3.
func NewConfigurationStorage(config model.ConfigStorageSettings) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from the S3 bucket...")

	s3client, err := s3Storage.NewS3Client(config.S3.Region)
	if err != nil {
		log.Fatalf("Cannot initialize S3 client: %s.", err)
	}

	cs := &ConfigurationStorage{
		Client:     s3client,
		Bucket:     config.S3.Bucket,
		ObjectName: config.S3.Key,
		UpdateChan: make(chan interface{}, 1),
		config:     config,
	}

	settings, err := cs.LoadServerSettings(true)
	if err != nil {
		return nil, fmt.Errorf("Cannot not load settings from etcd config storage: %s", err)
	}

	var keyStorage model.KeyStorage

	switch settings.KeyStorage.Type {
	case model.KeyStorageTypeLocal:
		keyStorage, err = keyStorageLocal.NewKeyStorage(settings.KeyStorage)
	case model.KeyStorageTypeS3:
		keyStorage, err = keyStorageS3.NewKeyStorage(settings.KeyStorage)
	default:
		return nil, fmt.Errorf("Unknown key storage type: %s", settings.KeyStorage.Type)
	}
	if err != nil {
		return nil, err
	}
	cs.keyStorage = keyStorage
	return cs, nil
}

// LoadServerSettings loads server configuration from S3.
func (cs *ConfigurationStorage) LoadServerSettings(forceReload bool) (model.ServerSettings, error) {
	if !forceReload && cs.cached {
		return cs.cache, nil
	}

	getObjInput := &s3.GetObjectInput{
		Bucket: aws.String(cs.Bucket),
		Key:    aws.String(cs.ObjectName),
	}

	resp, err := cs.Client.GetObject(getObjInput)
	if err != nil {
		return model.ServerSettings{}, fmt.Errorf("Cannot get object from S3: %s", err)
	}
	defer resp.Body.Close()

	var settings model.ServerSettings
	if err = yaml.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return model.ServerSettings{}, fmt.Errorf("Cannot decode S3 response: %s", err)
	}

	if settings.Config.Type != model.ConfigStorageTypeS3 {
		return model.ServerSettings{}, fmt.Errorf("Configuration file from S3 specifies configuration type %s", settings.Config.Type)
	}

	cs.cache = settings
	cs.cached = true

	return settings, nil
}

// WriteConfig puts new configuration into the storage.
func (cs *ConfigurationStorage) WriteConfig(settings model.ServerSettings) error {
	log.Println("Putting new config to S3...")

	valueBytes, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("Cannot marshal settings value: %s", err)
	}

	_, err = cs.Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(cs.Bucket),
		Key:           aws.String(cs.ObjectName),
		ACL:           aws.String("private"),
		StorageClass:  aws.String(s3.ObjectStorageClassStandard),
		Body:          bytes.NewReader(valueBytes),
		ContentLength: aws.Int64(int64(len(valueBytes))),
		ContentType:   aws.String("application/x-yaml"),
	})

	if err == nil {
		log.Println("Successfully put new configuration to S3")
	}

	// Indicate config update. To prevent writing to a closed channel, make a check.
	go func() {
		if cs.updateChanClosed {
			log.Println("Attempted to write to closed UpdateChan")
			return
		}
		cs.UpdateChan <- struct{}{}
	}()

	return err
}

// InsertKeys inserts new public and private keys to the S3 bucket.
func (cs *ConfigurationStorage) InsertKeys(keys *model.JWTKeys) error {
	if err := cs.keyStorage.InsertKeys(keys); err != nil {
		return err
	}
	return nil
}

// LoadKeys loads public and private keys from the key storage.
func (cs *ConfigurationStorage) LoadKeys(alg model.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
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
