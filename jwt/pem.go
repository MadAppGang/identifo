package jwt

import (
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

//LoadPrivateKeyFromPEM loads private key from PEM
func LoadPrivateKeyFromPEM(file string, alg model.TokenServiceAlgorithm) (interface{}, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var privateKey interface{}
	switch alg {
	case model.TokenServiceAlgorithmES256:
		privateKey, err = jwt.ParseECPrivateKeyFromPEM(prkb)
	case model.TokenServiceAlgorithmRS256:
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(prkb)
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	if err != nil {
		return nil, err
	}
	return privateKey, nil

}

//LoadPublicKeyFromPEM loads public key from PEM
func LoadPublicKeyFromPEM(file string, alg model.TokenServiceAlgorithm) (interface{}, error) {
	//load public key form pem file
	pkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var publicKey interface{}
	switch alg {
	case model.TokenServiceAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM(pkb)
	case model.TokenServiceAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pkb)
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

//LoadPublicKeyFromString loads public key from string
func LoadPublicKeyFromString(s string, alg model.TokenServiceAlgorithm) (interface{}, error) {
	var publicKey interface{}
	var err error
	switch alg {
	case model.TokenServiceAlgorithmES256:
		publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(s))
	case model.TokenServiceAlgorithmRS256:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(s))
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
