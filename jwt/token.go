package jwt

import (
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

var (
	ErrEmptyToken   = errors.New("Token is empty")
	ErrTokenInvalid = errors.New("Token is invalid")
)

//NewTokenService returns new JWT token service
//private is path to private key in pem format, please keep it in secret place
//public is path to the public key
//now we support only EC256 keypairs
func NewTokenService(private, public string) (model.TokenService, error) {
	t := TokenService{}

	//load private key from pem file
	prkb, err := ioutil.ReadFile(private)
	if err != nil {
		return nil, err
	}
	t.privateKey, err = jwt.ParseECPrivateKeyFromPEM(prkb)
	if err != nil {
		return nil, err
	}

	//load public key form pem file
	pkb, err := ioutil.ReadFile(public)
	if err != nil {
		return nil, err
	}
	t.publicKey, err = jwt.ParseECPublicKeyFromPEM(pkb)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

//TokenService JWT token service
type TokenService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

//Parse parses token data from string representation
func (ts *TokenService) Parse(s string) (model.Token, error) {
	tokenString := strings.TrimSpace(s)

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return ts.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	t := Token{}
	t.JWT = token
	return &t, nil
}

//NewToken creates new token for user
func (ts *TokenService) NewToken(u model.User, scopes []string) (model.Token, error) {
	//TODO: implementation
	return nil, nil
}

func (ts *TokenService) String(t model.Token) (string, error) {
	token, ok := t.(*Token)
	if !ok {
		return "", ErrTokenInvalid
	}
	if err := t.Validate(); err != nil {
		return "", err
	}
	if !token.JWT.Valid {
		return "", ErrTokenInvalid
	}
	str, err := token.JWT.SignedString(ts.privateKey)
	if err != nil {
		return "", err
	}

	return str, nil
}

//Token represents JWT token in the system
type Token struct {
	JWT *jwt.Token
}

//Validate validates token data, returns nil if all data is valid
func (t *Token) Validate() error {
	if t.JWT == nil {
		return ErrEmptyToken
	}
	if !t.JWT.Valid {
		return ErrTokenInvalid
	}

	return nil
}

//Claims extended claims structure
type Claims struct {
	Foo string `json:"user_profile,omitempty"`
	jwt.StandardClaims
}

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
