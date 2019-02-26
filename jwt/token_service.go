package jwt

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

var (
	ErrEmptyToken              = errors.New("Token is empty")
	ErrWrongSignatureAlgorithm = errors.New("Unsupported signature algorithm")
	ErrTokenInvalid            = errors.New("Token is invalid")
	ErrCreatingToken           = errors.New("Error creating token")
	ErrSavingToken             = errors.New("Error saving token")
	ErrInvalidApp              = errors.New("Application is not eligible to obtain the token")
	ErrInvalidOfflineScope     = errors.New("Requested scope don't have offline value")
	ErrInvalidUser             = errors.New("The user could not obtain the new token")

	//TokenLifespan expiry token time, one week
	TokenLifespan = int64(604800)
	//RefreshTokenLifespan default expire time for refresh token, one year
	RefreshTokenLifespan = int64(31557600)
)

const (
	PayloadName = "name"
)

//NewTokenService returns new JWT token service
//private is path to private key in pem format, please keep it in secret place
//public is path to the public key
//now we support only ES256 and RS256 keypairs
func NewTokenService(private, public, issuer string, alg model.TokenServiceAlgorithm, storage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage, options ...func(*TokenService) error) (model.TokenService, error) {
	t := TokenService{}
	t.issuer = issuer
	t.appStorage = appStorage
	t.userStorage = userStorage
	t.tokenStorage = storage
	// 2 hours is default expire time for refresh token
	t.resetTokenLifespan = int64(7200)
	// 2 days is default expire time for auth token
	t.authTokenLifespan = int64(60 * 60 * 24 * 2)

	//trying to guess algo from the private file
	if alg == model.TokenServiceAlgorithmAuto {
		_, err := LoadPrivateKeyFromPEM(private, model.TokenServiceAlgorithmES256)
		if err == nil {
			alg = model.TokenServiceAlgorithmES256
		} else {
			_, err := LoadPrivateKeyFromPEM(private, model.TokenServiceAlgorithmRS256)
			if err == nil {
				alg = model.TokenServiceAlgorithmRS256
			}
		}
		if alg == model.TokenServiceAlgorithmAuto {
			return nil, ErrWrongSignatureAlgorithm
		}
	}

	t.algorithm = alg
	//load private key from pem file
	var err error
	t.privateKey, err = LoadPrivateKeyFromPEM(private, alg)
	if err != nil {
		return nil, err
	}
	t.publicKey, err = LoadPublicKeyFromPEM(public, alg)
	if err != nil {
		return nil, err
	}

	// Apply options
	for _, option := range options {
		if err := option(&t); err != nil {
			return nil, err
		}
	}

	return &t, nil
}

//TokenService JWT token service
type TokenService struct {
	privateKey         interface{} //*ecdsa.PrivateKey, or *rsa.PrivateKey
	publicKey          interface{} //*ecdsa.PublicKey, or *rsa.PublicKey
	tokenStorage       model.TokenStorage
	appStorage         model.AppStorage
	userStorage        model.UserStorage
	algorithm          model.TokenServiceAlgorithm
	issuer             string
	resetTokenLifespan int64
	authTokenLifespan  int64
}

//Issuer returns issuer name
func (ts *TokenService) Issuer() string {
	return ts.issuer
}

//Algorithm - return signature algorithm
func (ts *TokenService) Algorithm() string {
	switch ts.algorithm {
	case model.TokenServiceAlgorithmES256:
		return "ES256"
	case model.TokenServiceAlgorithmRS256:
		return "RS256"
	default:
		return ""
	}
}

//PublicKey returns public key
func (ts *TokenService) PublicKey() interface{} {
	return ts.publicKey
}

//KeyID returns public key ID, using SHA-1 fingertip
func (ts *TokenService) KeyID() string {
	if der, err := x509.MarshalPKIXPublicKey(ts.publicKey); err == nil {
		s := sha1.Sum(der)
		return base64.RawURLEncoding.EncodeToString(s[:]) //slice from [20]byte
	}
	return ""
}

// AuthTokenLifespan return auth token lifespan
func (ts *TokenService) AuthTokenLifespan() int64 {
	return ts.authTokenLifespan
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

	payload := make(map[string]string)
	if contains(app.TokenPayload(), PayloadName) {
		payload[PayloadName] = u.Name()
	}
	now := TimeFunc().Unix()

	lifespan := app.TokenLifespan()
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    model.AccessTokenType,
		KeyID:   ts.KeyID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  app.ID(),
			IssuedAt:  now,
		},
	}
	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenServiceAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenServiceAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	token := jwt.NewWithClaims(sm, claims)
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

	payload := make(map[string]string)
	if contains(app.TokenPayload(), PayloadName) {
		payload[PayloadName] = u.Name()
	}
	now := TimeFunc().Unix()

	lifespan := app.RefreshTokenLifespan()
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    model.RefrestTokenType,
		KeyID:   ts.KeyID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  app.ID(),
			IssuedAt:  now,
		},
	}
	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenServiceAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenServiceAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	token := jwt.NewWithClaims(sm, claims)
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

//NewResetToken creates new token for reset password
func (ts *TokenService) NewResetToken(userID string) (model.Token, error) {
	now := TimeFunc().Unix()

	lifespan := ts.resetTokenLifespan

	claims := Claims{
		Type:  model.ResetTokenType,
		KeyID: ts.KeyID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   userID,
			Audience:  "identifo",
			IssuedAt:  now,
		},
	}
	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenServiceAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenServiceAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	token := jwt.NewWithClaims(sm, claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &Token{JWT: token, new: true}, nil
}

// NewAuthToken creates auth token
func (ts *TokenService) NewAuthToken(u model.User) (model.Token, error) {
	//check user
	if !u.Active() {
		return nil, ErrInvalidUser
	}
	now := TimeFunc().Unix()
	lifespan := ts.resetTokenLifespan

	claims := Claims{
		Type:  model.AuthTokenType,
		KeyID: ts.KeyID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID(),
			Audience:  "identifo",
			IssuedAt:  now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenServiceAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenServiceAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ErrWrongSignatureAlgorithm
	}
	token := jwt.NewWithClaims(sm, claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &Token{JWT: token, new: true}, nil
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

// ResetTokenLifespan sets custom lifespan in seconds for the reset token
func ResetTokenLifespan(lifespan int64) func(*TokenService) error {
	return func(ts *TokenService) error {
		ts.resetTokenLifespan = lifespan
		return nil
	}
}

// AuthTokenLifespan sets custom lifespan in seconds for the auth token
func AuthTokenLifespan(lifespan int64) func(*TokenService) error {
	return func(ts *TokenService) error {
		ts.authTokenLifespan = lifespan
		return nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.TrimSpace(strings.ToLower(a)) == strings.TrimSpace(strings.ToLower(e)) {
			return true
		}
	}
	return false
}
