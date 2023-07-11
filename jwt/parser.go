package jwt

import (
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/madappgang/identifo/v2/model"
)

// Parse parses token data from the string representation.
func ParseTokenString(str string) (model.Token, error) {
	tokenString := strings.TrimSpace(str)
	parser := jwt.Parser{}

	token, _, err := parser.ParseUnverified(tokenString, &model.Claims{})
	if err != nil {
		return nil, err
	}

	return &model.JWToken{JWT: token}, nil
}
