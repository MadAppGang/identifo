package jwt

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

var (
	ErrEmptyToken    = errors.New("Token is empty")
	ErrTokenInvalid  = errors.New("Token is invalid")
	ErrCreatingToken = errors.New("Error creating token")

	//TokenTimelife expiry token time, one week
	TokenTimelife = int64(604800)
)

//NewTokenService returns new JWT token service
//private is path to private key in pem format, please keep it in secret place
//public is path to the public key
//now we support only EC256 keypairs
func NewTokenService(private, public, issuer string) (model.TokenService, error) {
	t := TokenService{}
	t.issuer = issuer
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
	issuer     string
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
func (ts *TokenService) NewToken(u model.User, scopes []string, appID string) (model.Token, error) {

	profileString := ""
	profileBytes, err := json.Marshal(u.Profile())
	if err != nil {
		return nil, ErrCreatingToken
	} else {
		profileString = string(profileBytes)
	}
	now := TimeFunc().Unix()

	claims := Claims{
		Scopes:      strings.Join(scopes, " "),
		UserProfile: profileString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + TokenTimelife),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  appID,
			IssuedAt:  now,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &Token{JWT: token, new: true}, nil
}

func (ts *TokenService) String(t model.Token) (string, error) {
	token, ok := t.(*Token)
	if !ok {
		log.Println("Here !!")
		return "", ErrTokenInvalid
	}
	if err := t.Validate(); err != nil {
		log.Println("Here3 !!")
		return "", err
	}
	if !token.new && !token.JWT.Valid {
		log.Println("Here 2!!")
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
	new bool
}

//Validate validates token data, returns nil if all data is valid
func (t *Token) Validate() error {
	if t.JWT == nil {
		return ErrEmptyToken
	}
	if !t.new && !t.JWT.Valid {
		return ErrTokenInvalid
	}

	return nil
}

//Claims extended claims structure
type Claims struct {
	UserProfile string `json:"user_profile,omitempty"`
	Scopes      string
	jwt.StandardClaims
}

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
