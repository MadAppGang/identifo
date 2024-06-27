package s3

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
)

// KeyStorage is a wrapper over private key files
type KeyStorage struct {
	logger         *slog.Logger
	client         *s3.S3
	bucket         string
	privateKeyPath string
}

// NewKeyStorage creates and returns new S3-backed key files storage.
func NewKeyStorage(
	logger *slog.Logger,
	settings model.FileStorageS3,
) (*KeyStorage, error) {
	s3Client, err := NewS3Client(settings.Region, settings.Endpoint)
	if err != nil {
		return nil, err
	}

	return &KeyStorage{
		logger:         logger,
		client:         s3Client,
		bucket:         settings.Bucket,
		privateKeyPath: settings.Key,
	}, nil
}

// ReplaceKey replaces private  key into S3 key storage
func (ks *KeyStorage) ReplaceKey(keyPEM []byte) error {
	ks.logger.Info("Putting new keys to S3...")

	if keyPEM == nil {
		return fmt.Errorf("cannot insert empty key")
	}

	_, err := ks.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(ks.bucket),
		Key:          aws.String(ks.privateKeyPath),
		ACL:          aws.String("private"),
		StorageClass: aws.String(s3.ObjectStorageClassStandard),
		Body:         bytes.NewReader(keyPEM),
		ContentType:  aws.String("application/x-pem-file"),
	})
	if err == nil {
		ks.logger.Info("Successfully put key to S3",
			"path", ks.privateKeyPath)
	}

	return nil
}

// LoadPrivateKey loads private key from the storage
func (ks *KeyStorage) LoadPrivateKey() (interface{}, error) {
	getKeyInput := &s3.GetObjectInput{
		Bucket: aws.String(ks.bucket),
		Key:    aws.String(ks.privateKeyPath),
	}

	resp, err := ks.client.GetObject(getKeyInput)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s from S3: %w", ks.privateKeyPath, err)
	}
	defer resp.Body.Close()

	keyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot decode S3 response: %w", err)
	}

	privateKey, _, err := jwt.LoadPrivateKeyFromPEMString(string(keyData))
	if err != nil {
		return nil, fmt.Errorf("cannot load private key: %s", err)
	}

	return privateKey, nil
}
