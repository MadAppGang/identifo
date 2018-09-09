package jwt

import jwt "github.com/dgrijalva/jwt-go"

//Token represents JWT token in the system
type Token struct {
}

//Parse parses token from string
func Parse(t string) (*Token, error) {
	return nil, nil
}

//Claims extended claims structure
type Claims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}
