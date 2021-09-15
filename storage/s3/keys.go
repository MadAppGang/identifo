package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

// KeyStorage is a wrapper over private key files
type KeyStorage struct {
	Client         *s3.S3
	Bucket         string
	PrivateKeyPath string
}

// NewKeyStorage creates and returns new S3-backed key files storage.
func NewKeyStorage(settings model.S3KeyStorageSettings) (*KeyStorage, error) {
	s3Client, err := NewS3Client(settings.Region)
	if err != nil {
		return nil, err
	}

	return &KeyStorage{
		Client:         s3Client,
		Bucket:         settings.Bucket,
		PrivateKeyPath: settings.PrivateKeyKey,
	}, nil
}

// ReplaceKey replaces private  key into S3 key storage
func (ks *KeyStorage) ReplaceKey(keyPEM []byte) error {
	log.Println("Putting new keys to S3...")

	if keyPEM == nil {
		return fmt.Errorf("Cannot insert empty key")
	}

	_, err := ks.Client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(ks.Bucket),
		Key:          aws.String(ks.PrivateKeyPath),
		ACL:          aws.String("private"),
		StorageClass: aws.String(s3.ObjectStorageClassStandard),
		Body:         bytes.NewReader(keyPEM),
		ContentType:  aws.String("application/x-pem-file"),
	})
	if err == nil {
		log.Printf("Successfully put %s to S3\n", ks.PrivateKeyPath)
	}

	return nil
}

// LoadPrivateKey loads private key from the storage
func (ks *KeyStorage) LoadPrivateKey() (interface{}, error) {
	getKeyInput := &s3.GetObjectInput{
		Bucket: aws.String(ks.Bucket),
		Key:    aws.String(ks.PrivateKeyPath),
	}

	resp, err := ks.Client.GetObject(getKeyInput)
	if err != nil {
		return nil, fmt.Errorf("Cannot get %s from S3: %s", ks.PrivateKeyPath, err)
	}
	defer resp.Body.Close()

	keyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode S3 response: %s", err)
	}

	privateKey, _, err := jwt.LoadPrivateKeyFromPEMString(string(keyData))
	if err != nil {
		return nil, fmt.Errorf("cannot load private key: %s", err)
	}

	return privateKey, nil
}
