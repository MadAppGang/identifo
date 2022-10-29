package s3

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	"gopkg.in/yaml.v2"
)

// ConfigurationStorage is a server configuration storage in S3.
type ConfigurationStorage struct {
	Client           *s3.S3
	Bucket           string
	ObjectName       string
	UpdateChan       chan interface{}
	updateChanClosed bool
	config           model.FileStorageSettings
	cache            *model.ServerSettings
	errors           []error
}

// NewConfigurationStorage creates new server config storage in S3.
func NewConfigurationStorage(config model.FileStorageSettings) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from the S3 bucket...")

	if config.Type != model.FileStorageTypeS3 {
		return nil, fmt.Errorf("Configuration file from S3 specifies configuration type %s", config.Type)
	}

	s3client, err := NewS3Client(config.S3.Region, config.S3.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize S3 client: %s.", err)
	}

	cs := &ConfigurationStorage{
		Client:     s3client,
		Bucket:     config.S3.Bucket,
		ObjectName: config.S3.Key,
		UpdateChan: make(chan interface{}, 1),
		config:     config,
	}

	return cs, nil
}

// LoadServerSettings loads server configuration from S3.
func (cs *ConfigurationStorage) LoadServerSettings(validate bool) (model.ServerSettings, []error) {
	cs.errors = nil

	cs.errors = nil
	getObjInput := &s3.GetObjectInput{
		Bucket: aws.String(cs.Bucket),
		Key:    aws.String(cs.ObjectName),
	}

	resp, err := cs.Client.GetObject(getObjInput)
	if err != nil {
		e := fmt.Errorf("Cannot get object from S3: %s", err)
		cs.errors = append(cs.errors, e)
		return model.ServerSettings{}, cs.errors
	}
	defer resp.Body.Close()

	var settings model.ServerSettings
	if err = yaml.NewDecoder(resp.Body).Decode(&settings); err != nil {
		cs.errors = append(cs.errors, fmt.Errorf("Cannot decode S3 response: %s", err))
		return model.ServerSettings{}, cs.errors
	}

	settings.Config = cs.config
	cs.cache = &settings
	if validate {
		cs.errors = settings.Validate(true)
	}

	return settings, cs.errors
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
	return true
}

func (cs *ConfigurationStorage) LoadedSettings() *model.ServerSettings {
	return cs.cache
}

func (cs *ConfigurationStorage) Errors() []error {
	return cs.errors
}
