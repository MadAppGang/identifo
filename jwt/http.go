package jwt

import (
	"crypto/ecdsa"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

//TokenHeaderKeyPrefix token prefix regarding RFCXXX
const TokenHeaderKeyPrefix = "BEARER "

//ExtractTokenFromBearerHeader extracts token from bearer token header value
func ExtractTokenFromBearerHeader(token string) []byte {
	token = strings.TrimSpace(token)
	if (len(token) <= len(TokenHeaderKeyPrefix)) ||
		(strings.ToUpper(token[0:len(TokenHeaderKeyPrefix)]) != TokenHeaderKeyPrefix) {
		return nil
	}

	token = token[len(TokenHeaderKeyPrefix):]
	return []byte(token)
}

//ParseTokenWithPublicKey parses token with public key provided
func ParseTokenWithPublicKey(t string, publicKey *ecdsa.PublicKey) (model.Token, error) {
	tokenString := strings.TrimSpace(t)

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	resultToken := Token{}
	resultToken.JWT = token
	return &resultToken, nil
}

//LoadPrivateKeyFromPEM loads private key from PEM
func LoadPrivateKeyFromPEM(file string) (*ecdsa.PrivateKey, error) {
	prkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(prkb)
	if err != nil {
		return nil, err
	}
	return privateKey, nil

}

//LoadPublicKeyFromPEM loads public key from PEM
func LoadPublicKeyFromPEM(file string) (*ecdsa.PublicKey, error) {
	//load public key form pem file
	pkb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(pkb)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
