package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

func NewTokenService() model.TokenService {
	t := TokenService{}
	return &t
}

//TokenService JWT token service
type TokenService struct {
}

func (ts *TokenService) Parse(string) (model.Token, error) {
	//TODO: implementation
	return nil, nil
}

func (ts *TokenService) NewToken(model.User) (model.Token, error) {
	//TODO: implementation
	return nil, nil
}

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
	jwtgo.StandardClaims
}

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
