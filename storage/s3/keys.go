package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/madappgang/identifo/model"
)

// KeyStorage is a wrapper over public and private key files.
type KeyStorage struct {
	Client         *s3.S3
	Bucket         string
	PublicKeyPath  string
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
		PublicKeyPath:  settings.PublicKeyKey,
	}, nil
}

// InsertKeys inserts private and public keys into S3 key storage.
func (ks *KeyStorage) ReplaceKeys(keys model.JWTKeys) error {
	keysMap := map[string]interface{}{
		ks.PrivateKeyPath: keys.Private,
		ks.PublicKeyPath:  keys.Public,
	}
	log.Println("Putting new keys to S3...")

	for name, file := range keysMap {
		reader, ok := file.(io.ReadSeeker)
		if !ok {
			return fmt.Errorf("%s cannot be read", name)
		}

		_, err := ks.Client.PutObject(&s3.PutObjectInput{
			Bucket:       aws.String(ks.Bucket),
			Key:          aws.String(name),
			ACL:          aws.String("private"),
			StorageClass: aws.String(s3.ObjectStorageClassStandard),
			Body:         reader,
			ContentType:  aws.String("application/x-pem-file"),
		})
		if err == nil {
			log.Printf("Successfully put %s to S3\n", name)
		}
	}
	return nil
}

// LoadKeys loads keys from the key storage.
func (ks *KeyStorage) LoadKeys(alg model.TokenSignatureAlgorithm) (model.JWTKeys, error) {
	keys := model.JWTKeys{}

	for _, keyPath := range [2]string{ks.PublicKeyPath, ks.PrivateKeyPath} {
		getKeyInput := &s3.GetObjectInput{
			Bucket: aws.String(ks.Bucket),
			Key:    aws.String(keyPath),
		}

		resp, err := ks.Client.GetObject(getKeyInput)
		if err != nil {
			return keys, fmt.Errorf("Cannot get %s from S3: %s", keyPath, err)
		}
		defer resp.Body.Close()

		key, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return keys, fmt.Errorf("Cannot decode S3 response: %s", err)
		}

		if strings.Contains(keyPath, ks.PublicKeyPath) {
			keys.Public = key
			keys.Algorithm, err = ks.guessTokenServiceAlgorithm(key)
			if err != nil {
				return keys, err
			}
		} else {
			keys.Private = key
		}
	}
	return keys, nil
}

func (ks *KeyStorage) guessTokenServiceAlgorithm(publicKey []byte) (interface{}, error) {
	_, errES := jwt.ParseECPublicKeyFromPEM(publicKey)
	if errES == nil {
		return model.TokenSignatureAlgorithmES256, nil
	}
	_, errRS := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if errRS == nil {
		return model.TokenSignatureAlgorithmRS256, nil
	}
	return nil, fmt.Errorf("Cannot guess token service algorithm. It's neither ES256 (%s), nor RS256 (%s)", errES, errRS)
}

func (ks *KeyStorage) GetKeys() (model.KeysPEM, error) {
	keys := model.KeysPEM{}

	publk, err := ks.redFile(ks.PublicKeyPath)
	if err != nil {
		return keys, fmt.Errorf("Cannot get %s from S3: %s", ks.PublicKeyPath, err)
	}
	keys.Public = string(publk)

	private, err := ks.redFile(ks.PrivateKeyPath)
	if err != nil {
		return keys, fmt.Errorf("Cannot get %s from S3: %s", ks.PrivateKeyPath, err)
	}
	keys.Private = string(private)

	return keys, nil
}

func (ks *KeyStorage) redFile(key string) ([]byte, error) {
	getKeyInput := &s3.GetObjectInput{
		Bucket: aws.String(ks.Bucket),
		Key:    aws.String(key),
	}

	resp, err := ks.Client.GetObject(getKeyInput)
	if err != nil {
		return nil, fmt.Errorf("Cannot get %s from S3: %s", key, err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
