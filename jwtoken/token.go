package jwtoken

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

//NewTokenService returns new JWT token service
func NewTokenService() model.TokenService {
	t := TokenService{}
	return &t
}

//TokenService JWT token service
type TokenService struct {
}

//Parse parses tojen data from string representation
func (ts *TokenService) Parse(string) (model.Token, error) {
	//TODO: implementation
	return nil, nil
}

//NewToken creates new token for user
func (ts *TokenService) NewToken(u model.User, scopes []string) (model.Token, error) {
	//TODO: implementation
	return nil, nil
}

//Token represents JWT token in the system
type Token struct {
}

//Validate validates token data, returns nil if all data is valid
func (t *Token) Validate() error {
	return nil
}

//String returns string representation of the token
func (t *Token) String() string {
	return ""
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

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
