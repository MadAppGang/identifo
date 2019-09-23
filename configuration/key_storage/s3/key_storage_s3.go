package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	jwt "github.com/dgrijalva/jwt-go"
	s3Storage "github.com/madappgang/identifo/external_services/storage/s3"
	ijwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

// KeyStorage is a wrapper over public and private key files.
type KeyStorage struct {
	Client         *s3.S3
	Bucket         string
	PublicKeyName  string
	PrivateKeyName string
}

// NewKeyStorage creates and returns new S3-backed key files storage.
func NewKeyStorage(settings model.KeyStorageSettings) (*KeyStorage, error) {
	s3Client, err := s3Storage.NewS3Client(settings.Region)
	if err != nil {
		return nil, err
	}

	return &KeyStorage{
		Client:         s3Client,
		Bucket:         settings.Bucket,
		PublicKeyName:  settings.PublicKey,
		PrivateKeyName: settings.PrivateKey,
	}, nil
}

// InsertKeys inserts private and public keys into S3 key storage.
func (ks *KeyStorage) InsertKeys(keys *model.JWTKeys) error {
	if keys == nil {
		return fmt.Errorf("Empty keys")
	}
	keysMap := map[string]interface{}{
		ks.PrivateKeyName: keys.Private,
		ks.PublicKeyName:  keys.Public,
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
func (ks *KeyStorage) LoadKeys(alg ijwt.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
	keys := new(model.JWTKeys)

	for _, keyName := range [2]string{ks.PublicKeyName, ks.PrivateKeyName} {
		getKeyInput := &s3.GetObjectInput{
			Bucket: aws.String(ks.Bucket),
			Key:    aws.String(keyName),
		}

		resp, err := ks.Client.GetObject(getKeyInput)
		if err != nil {
			return nil, fmt.Errorf("Cannot get %s from S3: %s", keyName, err)
		}
		defer resp.Body.Close()

		key, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Cannot decode S3 response: %s", err)
		}

		if keyName == ks.PublicKeyName {
			keys.Public = key
			keys.Algorithm, err = ks.guessTokenServiceAlgorithm(key)
			if err != nil {
				return nil, err
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
		return ijwt.TokenSignatureAlgorithmES256, nil
	}
	_, errRS := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if errRS == nil {
		return ijwt.TokenSignatureAlgorithmRS256, nil
	}
	return nil, fmt.Errorf("Cannot guess token service algorithm. It's neither ES256 (%s), nor RS256 (%s)", errES, errRS)
}
