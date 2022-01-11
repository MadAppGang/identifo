package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/madappgang/identifo/v2/model"
)

var ErrKeyMustBePEMEncoded = errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")

func LoadPrivateKeyFromPEMString(s string) (interface{}, model.TokenSignatureAlgorithm, error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode([]byte(s)); block == nil {
		return nil, model.TokenSignatureAlgorithmInvalid, ErrKeyMustBePEMEncoded
	}

	pp, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}

	switch private := pp.(type) {
	case *rsa.PrivateKey:
		if private.Size() != 256 {
			return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("rsa private key size is unsupported, expecting 256, got: %d", private.Size())
		}
		return private, model.TokenSignatureAlgorithmRS256, nil
	case *ecdsa.PrivateKey:
		// check curve bits size and type
		if private.Curve.Params().BitSize != 256 {
			return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("ecdsa private key bit size is unsupported, expecting 256, got: %d", private.Curve.Params().BitSize)
		}
		if private.Curve.Params().Name != "P-256" {
			return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("ecdsa private key curve name us unsupported, expecting curve P-256, got: %s", private.Curve.Params().Name)
		}
		return private, model.TokenSignatureAlgorithmES256, nil
	default:
		return nil, model.TokenSignatureAlgorithmInvalid, fmt.Errorf("could not load unsupported key type: %T\n", private)
	}
}

// LoadPrivateKeyFromPEM loads private key from PEM file.
func LoadPrivateKeyFromPEMFile(file string) (interface{}, model.TokenSignatureAlgorithm, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, model.TokenSignatureAlgorithmInvalid, err
	}
	return LoadPrivateKeyFromPEMString(string(prkb))
}

// LoadPublicKeyFromString loads public key from string.
func LoadPublicKeyFromString(s string) (interface{}, model.TokenSignatureAlgorithm, error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode([]byte(s)); block == nil {
		return nil, model.TokenSignatureAlgorithmInvalid, ErrKeyMustBePEMEncoded
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
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
	b64 := []byte(base64.StdEncoding.EncodeToString(pk))
	return fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s-----END PRIVATE KEY-----\n", Make64ColsString(b64)), nil
}

func MarshalPublicKeyToPEM(key interface{}) (string, error) {
	pk, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("error creating PEM: %v", err)
	}
	b64 := []byte(base64.StdEncoding.EncodeToString(pk))
	return fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s-----END PUBLIC KEY-----\n", Make64ColsString(b64)), nil
}

func Make64ColsString(slice []byte) string {
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

func GenerateNewPrivateKey(alg model.TokenSignatureAlgorithm) (interface{}, error) {
	switch alg {
	case model.TokenSignatureAlgorithmRS256:
		return rsa.GenerateKey(rand.Reader, 2048)
	case model.TokenSignatureAlgorithmES256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		return nil, fmt.Errorf("unable to generate new private key, unsupported algorithm: %s\n", alg)
	}
}
