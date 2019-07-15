package s3

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3Storage "github.com/madappgang/identifo/external_services/storage/s3"
	"github.com/madappgang/identifo/model"
	"gopkg.in/yaml.v2"
)

const (
	identifoConfigS3BucketEnvName = "IDENTIFO_CONFIG_BUCKET"
)

// ConfigurationStorage is a server configuration storage in S3.
type ConfigurationStorage struct {
	Client           *s3.S3
	Bucket           string
	Key              string
	UpdateChan       chan interface{}
	updateChanClosed bool
}

// NewConfigurationStorage creates new server config storage in S3.
func NewConfigurationStorage(settings model.ConfigurationStorageSettings) (*ConfigurationStorage, error) {
	s3client, err := s3Storage.NewS3Client(settings.Region)
	if err != nil {
		return nil, err
	}

	bucket := os.Getenv(identifoConfigS3BucketEnvName)
	if len(bucket) == 0 {
		return nil, fmt.Errorf("No %s specified", identifoConfigS3BucketEnvName)
	}

	if len(settings.SettingsKey) == 0 {
		return nil, fmt.Errorf("No file key for the bucket specified")
	}

	cs := &ConfigurationStorage{
		Client:     s3client,
		Bucket:     bucket,
		Key:        settings.SettingsKey,
		UpdateChan: make(chan interface{}, 1),
	}
	return cs, nil
}

// LoadServerSettings loads server configuration from S3.
func (cs *ConfigurationStorage) LoadServerSettings(settings *model.ServerSettings) error {
	getObjInput := &s3.GetObjectInput{
		Bucket: aws.String(cs.Bucket),
		Key:    aws.String(cs.Key),
	}

	resp, err := cs.Client.GetObject(getObjInput)
	if err != nil {
		return fmt.Errorf("Cannot get object from S3: %s", err)
	}
	defer resp.Body.Close()

	if err = yaml.NewDecoder(resp.Body).Decode(settings); err != nil {
		return fmt.Errorf("Cannot decode S3 response: %s", err)
	}

	if settings.ConfigurationStorage.Type != model.ConfigurationStorageTypeS3 {
		return fmt.Errorf("Configuration file from S3 specifies configuration type %s", settings.ConfigurationStorage.Type)
	}

	return nil
}

// Insert puts new configuration into the storage.
func (cs *ConfigurationStorage) Insert(key string, value interface{}) error {
	log.Println("Putting new config to S3...")

	valueBytes, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("Cannot marshal settings value: %s", err)
	}

	_, err = cs.Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(cs.Bucket),
		Key:           aws.String(cs.Key),
		ACL:           aws.String("private"),
		StorageClass:  aws.String("ONEZONE_IA"),
		Body:          bytes.NewReader(valueBytes),
		ContentLength: aws.Int64(int64(len(valueBytes))),
		ContentType:   aws.String("application/x-yaml"),
	})

	if err == nil {
		log.Println("Successfully put new configuration to S3")
	}

	// Indicate config update. To prevent writing to closed channel, make a check.
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
