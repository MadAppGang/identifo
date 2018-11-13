package jwt

import (
	"crypto/ecdsa"
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
