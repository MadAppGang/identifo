package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/madappgang/identifo/model"
)

func LoadPrivateKeyFromString(s string) (interface{}, model.TokenSignatureAlgorithm, error) {
	pp, err := x509.ParsePKCS8PrivateKey([]byte(s))
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}

	switch private := pp.(type) {
	case *rsa.PrivateKey:
		return private, model.TokenSignatureAlgorithmRS256, nil
	case *ecdsa.PrivateKey:
		return private, model.TokenSignatureAlgorithmES256, nil
	default:
		return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("could not load unsupported key type: %T\n", private)
	}
}

// LoadPrivateKeyFromPEM loads private key from PEM file.
func LoadPrivateKeyFromPEM(file string) (interface{}, model.TokenSignatureAlgorithm, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}
	return LoadPrivateKeyFromString(string(prkb))
}

// LoadPublicKeyFromString loads public key from string.
func LoadPublicKeyFromString(s string) (interface{}, model.TokenSignatureAlgorithm, error) {
	pub, err := x509.ParsePKIXPublicKey([]byte(s))
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, model.TokenSignatureAlgorithmRS256, nil
	case *ecdsa.PublicKey:
		return pub, model.TokenSignatureAlgorithmES256, nil
	default:
		return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("could not load unsupported key type: %T\n", pub)
	}
}

// LoadPublicKeyFromPEM loads public key from file
func LoadPublicKeyFromPEM(file string) (interface{}, model.TokenSignatureAlgorithm, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}

	return LoadPublicKeyFromString(string(prkb))
}

func MarshalPrivateKeyToPEM(key interface{}) (string, error) {
	pk, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("error creating PEM: %v", err)
	}
	b64 := []byte(base64.RawStdEncoding.EncodeToString(pk))
	return fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s-----END PRIVATE KEY-----\n", make64ColsString(b64)), nil
}

func MarshalPublicKeyToPEM(key interface{}) (string, error) {
	pk, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("error creating PEM: %v", err)
	}
	b64 := []byte(base64.RawStdEncoding.EncodeToString(pk))
	return fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s-----END PUBLIC KEY-----\n", make64ColsString(b64)), nil
}

func make64ColsString(slice []byte) string {
	chunks := chunkSlice(slice, 64)

	result := ""
	for _, line := range chunks {
		result = result + string(line) + "\n"
	}
	return result
}

// chunkSlice split slices
func chunkSlice(slice []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
