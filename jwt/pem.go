package jwt

import (
	"io/ioutil"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/madappgang/identifo/model"
)

var supportedSignatureAlgorithms = []model.TokenSignatureAlgorithm{model.TokenSignatureAlgorithmES256, model.TokenSignatureAlgorithmRS256}

// LoadPrivateKeyFromPEM loads private key from PEM file.
func LoadPrivateKeyFromPEM(file string, alg model.TokenSignatureAlgorithm) (interface{}, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var privateKey interface{}
	switch alg {
	case model.TokenSignatureAlgorithmES256:
		privateKey, err = jwt.ParseECPrivateKeyFromPEM(prkb)
	case model.TokenSignatureAlgorithmRS256:
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(prkb)
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// LoadPublicKeyFromPEM loads public key from PEM file.
func LoadPublicKeyFromPEM(file string, alg model.TokenSignatureAlgorithm) (interface{}, error) {
	if alg == model.TokenSignatureAlgorithmAuto {
		k, _, e := LoadPublicKeyFromPEMAuto(file)
		return k, e
	}

	pkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var publicKey interface{}
	switch alg {
	case model.TokenSignatureAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM(pkb)
	case model.TokenSignatureAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pkb)
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// LoadPublicKeyFromPEMAuto loads keys from pem file with key algorithm auto detection
func LoadPublicKeyFromPEMAuto(file string) (interface{}, model.TokenSignatureAlgorithm, error) {
	var err error
	var key interface{}
	alg := model.TokenSignatureAlgorithmAuto
	for _, a := range supportedSignatureAlgorithms {
		if key, err = LoadPublicKeyFromPEM(file, a); err == nil {
			alg = a
			break
		}
	}
	return key, alg, err
}

// LoadPublicKeyFromString loads public key from string.
func LoadPublicKeyFromString(s string, alg model.TokenSignatureAlgorithm) (interface{}, error) {
	if alg == model.TokenSignatureAlgorithmAuto {
		k, _, e := LoadPublicKeyFromStringAuto(s)
		return k, e
	}

	var publicKey interface{}
	var err error

	switch alg {
	case model.TokenSignatureAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(s))
	case model.TokenSignatureAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(s))
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// LoadPublicKeyFromStringAuto loads keys from string with key algorithm auto detection
func LoadPublicKeyFromStringAuto(s string) (interface{}, model.TokenSignatureAlgorithm, error) {
	var err error
	var key interface{}
	alg := model.TokenSignatureAlgorithmAuto
	for _, a := range supportedSignatureAlgorithms {
		if key, err = LoadPublicKeyFromString(s, a); err == nil {
			alg = a
			break
		}
	}
	return key, alg, err
}
