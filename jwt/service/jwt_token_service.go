package service

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	ijwt "github.com/madappgang/identifo/jwt"
	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
)

var (
	// ErrKeyFileNotFound is when key file not found.
	ErrKeyFileNotFound = errors.New("Key file not found")
	// ErrCreatingToken is a token creation error.
	ErrCreatingToken = errors.New("Error creating token")
	// ErrSavingToken is a token saving error.
	ErrSavingToken = errors.New("Error saving token")
	// ErrInvalidApp is when the application is not eligible to obtain the token
	ErrInvalidApp = errors.New("Application is not eligible to obtain the token")
	// ErrInvalidOfflineScope is when the requested scope does not have an offline value.
	ErrInvalidOfflineScope = errors.New("Requested scope don't have offline value")
	// ErrInvalidUser is when the user cannot obtain the new token.
	ErrInvalidUser = errors.New("The user cannot obtain the new token")

	// TokenLifespan is a token expiration time, one week.
	TokenLifespan = int64(604800) // int64(1*7*24*60*60)
	// InviteTokenLifespan is an invite token expiration time, one hour.
	InviteTokenLifespan = int64(3600) // int64(1*60*60)
	// RefreshTokenLifespan is a default expiration time for refresh tokens, one year.
	RefreshTokenLifespan = int64(31536000) // int(365*24*60*60)
)

const (
	// PayloadName is a JWT token payload name.
	PayloadName = "name"
)

// NewJWTokenService returns new JWT token service.
// Arguments:
// - privateKeyPath - the path to the private key in pem format. Please keep it in a secret place.
// - publicKeyPath - the path to the public key.
func NewJWTokenService(privateKeyPath, publicKeyPath, issuer string, alg ijwt.TokenSignatureAlgorithm, tokenStorage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage, options ...func(TokenService) error) (TokenService, error) {
	if _, err := os.Stat(privateKeyPath); err != nil {
		return nil, ErrKeyFileNotFound
	}
	if _, err := os.Stat(publicKeyPath); err != nil {
		return nil, ErrKeyFileNotFound
	}

	var privateKey interface{}
	var err error

	// Trying to guess algo from the private key file.
	if alg == ijwt.TokenSignatureAlgorithmAuto {
		if privateKey, err = ijwt.LoadPrivateKeyFromPEM(privateKeyPath, ijwt.TokenSignatureAlgorithmES256); err == nil {
			alg = ijwt.TokenSignatureAlgorithmES256
		} else if privateKey, err = ijwt.LoadPrivateKeyFromPEM(privateKeyPath, ijwt.TokenSignatureAlgorithmRS256); err == nil {
			alg = ijwt.TokenSignatureAlgorithmRS256
		}
	}
	if alg == ijwt.TokenSignatureAlgorithmAuto {
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	publicKey, err := ijwt.LoadPublicKeyFromPEM(publicKeyPath, alg)
	if err != nil {
		return nil, err
	}

	t := &JWTokenService{
		issuer:                 issuer,
		tokenStorage:           tokenStorage,
		appStorage:             appStorage,
		userStorage:            userStorage,
		resetTokenLifespan:     int64(2 * 60 * 60),      // 2 hours is a default expiration time for refresh tokens.
		webCookieTokenLifespan: int64(2 * 24 * 60 * 60), // 2 days is a default default expiration time for access tokens.
		algorithm:              alg,
		privateKey:             privateKey,
		publicKey:              publicKey,
	}

	// Apply options.
	for _, option := range options {
		if err := option(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// JWTokenService is a JWT token service.
type JWTokenService struct {
	privateKey             interface{} // *ecdsa.PrivateKey, or *rsa.PrivateKey
	publicKey              interface{} // *ecdsa.PublicKey, or *rsa.PublicKey
	tokenStorage           model.TokenStorage
	appStorage             model.AppStorage
	userStorage            model.UserStorage
	algorithm              ijwt.TokenSignatureAlgorithm
	issuer                 string
	resetTokenLifespan     int64
	webCookieTokenLifespan int64
}

// Issuer returns token issuer name.
func (ts *JWTokenService) Issuer() string {
	return ts.issuer
}

// Algorithm  returns signature algorithm.
func (ts *JWTokenService) Algorithm() string {
	switch ts.algorithm {
	case ijwt.TokenSignatureAlgorithmES256:
		return "ES256"
	case ijwt.TokenSignatureAlgorithmRS256:
		return "RS256"
	default:
		return ""
	}
}

// PublicKey returns public key.
func (ts *JWTokenService) PublicKey() interface{} {
	return ts.publicKey
}

// KeyID returns public key ID, using SHA-1 fingerprint.
func (ts *JWTokenService) KeyID() string {
	if der, err := x509.MarshalPKIXPublicKey(ts.publicKey); err == nil {
		s := sha1.Sum(der)
		return base64.RawURLEncoding.EncodeToString(s[:]) //slice from [20]byte
	}
	return ""
}

// WebCookieTokenLifespan return auth token lifespan
func (ts *JWTokenService) WebCookieTokenLifespan() int64 {
	return ts.webCookieTokenLifespan
}

// Parse parses token data from the string representation.
func (ts *JWTokenService) Parse(s string) (ijwt.Token, error) {
	tokenString := strings.TrimSpace(s)

	token, err := jwt.ParseWithClaims(tokenString, &ijwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counterpart to verify them.
		return ts.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return &ijwt.JWToken{JWT: token}, nil
}

// ValidateTokenString parses token and validates it.
func (ts *JWTokenService) ValidateTokenString(tstr string, v jwtValidator.Validator, tokenType string) (ijwt.Token, error) {
	token, err := ts.Parse(tstr)
	if err != nil {
		return nil, err
	}

	if err := v.Validate(token); err != nil {
		return nil, err
	}

	if token.Type() != tokenType {
		return nil, err
	}

	return token, nil
}

// NewToken creates new token for user.
func (ts *JWTokenService) NewToken(u model.User, scopes []string, app model.AppData) (ijwt.Token, error) {
	if !app.Active() {
		return nil, ErrInvalidApp
	}

	if !u.Active() {
		return nil, ErrInvalidUser
	}

	payload := make(map[string]string)
	if contains(app.TokenPayload(), PayloadName) {
		payload[PayloadName] = u.Username()
	}
	now := ijwt.TimeFunc().Unix()

	lifespan := app.TokenLifespan()
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := ijwt.Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    AccessTokenType,
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
	case ijwt.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case ijwt.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	token := ijwt.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &ijwt.JWToken{JWT: token, New: true}, nil
}

// NewInviteToken creates new invite token.
func (ts *JWTokenService) NewInviteToken() (ijwt.Token, error) {
	payload := make(map[string]string)
	// add payload data here

	now := ijwt.TimeFunc().Unix()

	lifespan := InviteTokenLifespan

	claims := ijwt.Claims{
		Payload: payload,
		Type:    InviteTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + lifespan,
			Issuer:    ts.issuer,
			// Subject:   u.ID(),
			Audience: "identifo",
			IssuedAt: now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case ijwt.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case ijwt.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	token := ijwt.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &ijwt.JWToken{JWT: token, New: true}, nil
}

// NewRefreshToken creates new refresh token.
func (ts *JWTokenService) NewRefreshToken(u model.User, scopes []string, app model.AppData) (ijwt.Token, error) {
	if !app.Active() || !app.Offline() {
		return nil, ErrInvalidApp

	}
	// no offline request
	if !contains(scopes, OfflineScope) {
		return nil, ErrInvalidOfflineScope
	}

	if !u.Active() {
		return nil, ErrInvalidUser
	}

	payload := make(map[string]string)
	if contains(app.TokenPayload(), PayloadName) {
		payload[PayloadName] = u.Username()
	}
	now := ijwt.TimeFunc().Unix()

	lifespan := app.RefreshTokenLifespan()
	if lifespan == 0 {
		lifespan = RefreshTokenLifespan
	}

	claims := ijwt.Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    RefrestTokenType,
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
	case ijwt.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case ijwt.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	token := ijwt.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	t := &ijwt.JWToken{JWT: token, New: true}
	tokenString, err := ts.String(t)
	if err != nil {
		return nil, ErrSavingToken
	}

	if err := ts.tokenStorage.SaveToken(tokenString); err != nil {
		return nil, ErrSavingToken
	}
	return t, nil
}

// RefreshToken issues the new access token with access token
func (ts *JWTokenService) RefreshToken(refreshToken ijwt.Token) (ijwt.Token, error) {
	rt, ok := refreshToken.(*ijwt.JWToken)
	if !ok || rt == nil {
		return nil, ijwt.ErrTokenInvalid
	}

	if err := rt.Validate(); err != nil {
		return nil, err
	}

	claims, ok := rt.JWT.Claims.(*ijwt.Claims)
	if !ok || claims == nil {
		return nil, ijwt.ErrTokenInvalid
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

// NewResetToken creates new token for password resetting.
func (ts *JWTokenService) NewResetToken(userID string) (ijwt.Token, error) {
	now := ijwt.TimeFunc().Unix()

	lifespan := ts.resetTokenLifespan

	claims := ijwt.Claims{
		Type: ResetTokenType,
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
	case ijwt.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case ijwt.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	token := ijwt.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &ijwt.JWToken{JWT: token, New: true}, nil
}

// NewWebCookieToken creates new web cookie token.
func (ts *JWTokenService) NewWebCookieToken(u model.User) (ijwt.Token, error) {
	if !u.Active() {
		return nil, ErrInvalidUser
	}
	now := ijwt.TimeFunc().Unix()
	lifespan := ts.resetTokenLifespan

	claims := ijwt.Claims{
		Type: WebCookieTokenType,
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
	case ijwt.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case ijwt.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, ijwt.ErrWrongSignatureAlgorithm
	}

	token := ijwt.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &ijwt.JWToken{JWT: token, New: true}, nil
}

// String returns string representation of a token.
func (ts *JWTokenService) String(t ijwt.Token) (string, error) {
	token, ok := t.(*ijwt.JWToken)
	if !ok {
		return "", ijwt.ErrTokenInvalid
	}

	if err := t.Validate(); err != nil {
		return "", err
	}
	if !token.New && !token.JWT.Valid {
		return "", ijwt.ErrTokenInvalid
	}

	str, err := token.JWT.SignedString(ts.privateKey)
	if err != nil {
		return "", err
	}
	return str, nil
}

// ResetTokenLifespan sets custom lifespan in seconds for the reset token
func ResetTokenLifespan(lifespan int64) func(*JWTokenService) error {
	return func(ts *JWTokenService) error {
		ts.resetTokenLifespan = lifespan
		return nil
	}
}

// WebCookieTokenLifespan sets custom lifespan in seconds for the web cookie token
func WebCookieTokenLifespan(lifespan int64) func(*JWTokenService) error {
	return func(ts *JWTokenService) error {
		ts.webCookieTokenLifespan = lifespan
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
