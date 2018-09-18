package jwt

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

var (
	ErrEmptyToken          = errors.New("Token is empty")
	ErrTokenInvalid        = errors.New("Token is invalid")
	ErrCreatingToken       = errors.New("Error creating token")
	ErrSavingToken         = errors.New("Error saving token")
	ErrInvalidApp          = errors.New("Application is not eligible to obtain the token")
	ErrInvalidOfflineScope = errors.New("Requested scope don't have offline value")
	ErrInvalidUser         = errors.New("The user could not obtain the new token")

	//TokenLifespan expiry token time, one week
	TokenLifespan = int64(604800)
	//RefreshTokenLifespan default expire time for refresh token, one year
	RefreshTokenLifespan = int64(31557600)
)

//NewTokenService returns new JWT token service
//private is path to private key in pem format, please keep it in secret place
//public is path to the public key
//now we support only EC256 keypairs
func NewTokenService(private, public, issuer string, storage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage) (model.TokenService, error) {
	t := TokenService{}
	t.issuer = issuer
	t.appStorage = appStorage
	t.userStorage = userStorage
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
	privateKey   *ecdsa.PrivateKey
	publicKey    *ecdsa.PublicKey
	tokenStorage model.TokenStorage
	appStorage   model.AppStorage
	userStorage  model.UserStorage
	issuer       string
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
func (ts *TokenService) NewToken(u model.User, scopes []string, app model.AppData) (model.Token, error) {
	if !app.Active() {
		return nil, ErrInvalidApp
	}
	//check user
	if !u.Active() {
		return nil, ErrInvalidUser
	}

	profileString := ""
	profileBytes, err := json.Marshal(u.Profile())
	if err != nil {
		return nil, ErrCreatingToken
	}
	profileString = string(profileBytes)
	now := TimeFunc().Unix()

	lifespan := app.TokenLifespan()
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := Claims{
		Scopes:      strings.Join(scopes, " "),
		UserProfile: profileString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  app.ID(),
			IssuedAt:  now,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &Token{JWT: token, new: true}, nil
}

//NewRefreshToken creates new refresh token for the user
func (ts *TokenService) NewRefreshToken(u model.User, scopes []string, app model.AppData) (model.Token, error) {
	if !app.Active() || !app.Offline() {
		return nil, ErrInvalidApp

	}
	//no offline request
	if !contains(scopes, model.OfflineScope) {
		return nil, ErrInvalidOfflineScope
	}
	//check user
	if !u.Active() {
		return nil, ErrInvalidUser
	}
	profileString := ""
	profileBytes, err := json.Marshal(u.Profile())
	if err != nil {
		return nil, ErrCreatingToken
	}
	profileString = string(profileBytes)
	now := TimeFunc().Unix()

	lifespan := app.RefreshTokenLifespan()
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := Claims{
		Scopes:      strings.Join(scopes, " "),
		UserProfile: profileString,
		Type:        model.RefrestTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  app.ID(),
			IssuedAt:  now,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	t := Token{JWT: token, new: true}
	tokenString, err := ts.String(&t)
	if err != nil {
		return nil, ErrSavingToken
	}
	if err := ts.tokenStorage.SaveToken(tokenString); err != nil {
		return nil, ErrSavingToken
	}
	return &t, nil
}

//RefreshToken issues the new access token with access token
func (ts *TokenService) RefreshToken(refreshToken model.Token) (model.Token, error) {
	rt, ok := refreshToken.(*Token)
	if !ok || rt == nil {
		return nil, ErrTokenInvalid
	}
	if err := rt.Validate(); err != nil {
		return nil, err
	}
	claims, ok := rt.JWT.Claims.(*Claims)
	if !ok || claims == nil {
		return nil, ErrTokenInvalid
	}
	app, err := ts.appStorage.AppByID(claims.Audience)
	if err != nil || app == nil || !app.Offline() {
		return nil, ErrInvalidApp
	}

	user, err := ts.userStorage.UserByID(claims.Subject)
	if err != nil || user == nil || !user.Active() {
		return nil, ErrInvalidUser
	}
	token, err := ts.NewToken(user, strings.Split(claims.Scopes, " "), app)
	if err != nil {
		return nil, err
	}

	tokenString, err := ts.String(token)
	if err != nil {
		return nil, ErrSavingToken
	}
	if err := ts.tokenStorage.SaveToken(tokenString); err != nil {
		return nil, ErrSavingToken
	}
	return token, nil
}

func (ts *TokenService) String(t model.Token) (string, error) {
	token, ok := t.(*Token)
	if !ok {
		return "", ErrTokenInvalid
	}
	if err := t.Validate(); err != nil {
		return "", err
	}
	if !token.new && !token.JWT.Valid {
		return "", ErrTokenInvalid
	}
	str, err := token.JWT.SignedString(ts.privateKey)
	if err != nil {
		return "", err
	}

	return str, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.TrimSpace(strings.ToLower(a)) == strings.TrimSpace(strings.ToLower(e)) {
			return true
		}
	}
	return false
}
