package jwt

import (
	"strings"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/madappgang/identifo/model"
)

// TokenHeaderKeyPrefix is a token prefix regarding RFCXXX.
const TokenHeaderKeyPrefix = "BEARER "

// ExtractTokenFromBearerHeader extracts token from the Bearer token header value.
func ExtractTokenFromBearerHeader(token string) []byte {
	token = strings.TrimSpace(token)
	if (len(token) <= len(TokenHeaderKeyPrefix)) ||
		(strings.ToUpper(token[0:len(TokenHeaderKeyPrefix)]) != TokenHeaderKeyPrefix) {
		return nil
	}

	token = token[len(TokenHeaderKeyPrefix):]
	return []byte(token)
}

// ParseTokenWithPublicKey parses token with provided public key.
func ParseTokenWithPublicKey(t string, publicKey interface{}) (model.Token, error) {
	tokenString := strings.TrimSpace(t)

	parsedToken, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return &model.JWToken{JWT: parsedToken}, nil
}

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value. This is useful for testing or if your
// server uses a time zone different from your tokens'.
var TimeFunc = time.Now
