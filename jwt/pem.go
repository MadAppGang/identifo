package jwt

import (
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

// LoadPrivateKeyFromPEM loads private key from PEM file.
func LoadPrivateKeyFromPEM(file string, alg TokenSignatureAlgorithm) (interface{}, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var privateKey interface{}
	switch alg {
	case TokenSignatureAlgorithmES256:
		privateKey, err = jwt.ParseECPrivateKeyFromPEM(prkb)
	case TokenSignatureAlgorithmRS256:
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(prkb)
	default:
		return nil, ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return privateKey, nil

}

// LoadPublicKeyFromPEM loads public key from PEM file.
func LoadPublicKeyFromPEM(file string, alg TokenSignatureAlgorithm) (interface{}, error) {
	pkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var publicKey interface{}
	switch alg {
	case TokenSignatureAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM(pkb)
	case TokenSignatureAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pkb)
	default:
		return nil, ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// LoadPublicKeyFromString loads public key from string.
func LoadPublicKeyFromString(s string, alg TokenSignatureAlgorithm) (interface{}, error) {
	var publicKey interface{}
	var err error

	switch alg {
	case TokenSignatureAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(s))
	case TokenSignatureAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(s))
	default:
		return nil, ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
