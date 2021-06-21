package service

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/form3tech-oss/jwt-go"
	ijwt "github.com/madappgang/identifo/jwt"
	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
)

var (
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
	// PayloadName is a JWT token payload "name".
	PayloadName = "name"
)

// NewJWTokenService returns new JWT token service.
// Arguments:
// - privateKeyPath - the path to the private key in pem format. Please keep it in a secret place.
// - publicKeyPath - the path to the public key.
func NewJWTokenService(keys model.JWTKeys, issuer string, tokenStorage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage, options ...func(model.TokenService) error) (model.TokenService, error) {
	if keys.Private == nil || keys.Public == nil {
		return nil, fmt.Errorf("One of the keys is empty, or both")
	}

	tokenServiceAlg, ok := keys.Algorithm.(model.TokenSignatureAlgorithm)
	if !ok || tokenServiceAlg == model.TokenSignatureAlgorithmAuto {
		return nil, fmt.Errorf("Unknown token service algorithm %s ", keys.Algorithm)
	}

	t := &JWTokenService{
		issuer:                 issuer,
		tokenStorage:           tokenStorage,
		appStorage:             appStorage,
		userStorage:            userStorage,
		resetTokenLifespan:     int64(2 * 60 * 60),      // 2 hours is a default expiration time for refresh tokens.
		webCookieTokenLifespan: int64(2 * 24 * 60 * 60), // 2 days is a default default expiration time for access tokens.
		algorithm:              tokenServiceAlg,
		privateKey:             keys.Private,
		publicKey:              keys.Public,
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
	algorithm              model.TokenSignatureAlgorithm
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
	case model.TokenSignatureAlgorithmES256:
		return "ES256"
	case model.TokenSignatureAlgorithmRS256:
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
		return base64.RawURLEncoding.EncodeToString(s[:]) // slice from [20]byte
	}
	return ""
}

// WebCookieTokenLifespan return auth token lifespan
func (ts *JWTokenService) WebCookieTokenLifespan() int64 {
	return ts.webCookieTokenLifespan
}

// Parse parses token data from the string representation.
func (ts *JWTokenService) Parse(s string) (model.Token, error) {
	tokenString := strings.TrimSpace(s)

	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counterpart to verify them.
		return ts.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return &model.JWToken{JWT: token}, nil
}

// ValidateTokenString parses token and validates it.
func (ts *JWTokenService) ValidateTokenString(tstr string, v jwtValidator.Validator, tokenType string) (model.Token, error) {
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

// NewAccessToken creates new access token for user.
func (ts *JWTokenService) NewAccessToken(u model.User, scopes []string, app model.AppData, requireTFA bool, tokenPayload map[string]interface{}) (model.Token, error) {
	if !app.Active {
		return nil, ErrInvalidApp
	}

	if !u.Active {
		return nil, ErrInvalidUser
	}

	payload := make(map[string]interface{})
	if contains(app.TokenPayload, PayloadName) {
		payload[PayloadName] = u.Username
	}

	tokenType := model.TokenTypeAccess
	if requireTFA {
		scopes = []string{model.TokenTFAPreauthScope}
	}
	if len(tokenPayload) > 0 {
		for k, v := range tokenPayload {
			payload[k] = v
		}
	}

	now := ijwt.TimeFunc().Unix()

	lifespan := app.TokenLifespan
	if lifespan == 0 {
		lifespan = TokenLifespan
	}

	claims := model.Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID,
			Audience:  []string{app.ID},
			IssuedAt:  now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	token := model.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &model.JWToken{JWT: token, New: true}, nil
}

// NewRefreshToken creates new refresh token.
func (ts *JWTokenService) NewRefreshToken(u model.User, scopes []string, app model.AppData) (model.Token, error) {
	if !app.Active || !app.Offline {
		return nil, ErrInvalidApp
	}
	// no offline request
	if !contains(scopes, model.OfflineScope) {
		return nil, ErrInvalidOfflineScope
	}

	if !u.Active {
		return nil, ErrInvalidUser
	}

	payload := make(map[string]interface{})
	if contains(app.TokenPayload, PayloadName) {
		payload[PayloadName] = u.Username
	}
	now := ijwt.TimeFunc().Unix()

	lifespan := app.RefreshTokenLifespan
	if lifespan == 0 {
		lifespan = RefreshTokenLifespan
	}

	claims := model.Claims{
		Scopes:  strings.Join(scopes, " "),
		Payload: payload,
		Type:    model.TokenTypeRefresh,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID,
			Audience:  []string{app.ID},
			IssuedAt:  now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	token := model.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	t := &model.JWToken{JWT: token, New: true}
	tokenString, err := ts.String(t)
	if err != nil {
		return nil, ErrSavingToken
	}

	if err := ts.tokenStorage.SaveToken(tokenString); err != nil {
		return nil, ErrSavingToken
	}
	return t, nil
}

// RefreshAccessToken issues new access token for provided refresh token.
func (ts *JWTokenService) RefreshAccessToken(refreshToken model.Token) (model.Token, error) {
	rt, ok := refreshToken.(*model.JWToken)
	if !ok || rt == nil {
		return nil, model.ErrTokenInvalid
	}

	if err := rt.Validate(); err != nil {
		return nil, err
	}

	claims, ok := rt.JWT.Claims.(*model.Claims)
	if !ok || claims == nil {
		return nil, model.ErrTokenInvalid
	}

	app, err := ts.appStorage.AppByID(claims.Audience[0])
	if err != nil || !app.Offline {
		return nil, ErrInvalidApp
	}

	user, err := ts.userStorage.UserByID(claims.Subject)
	if err != nil || !user.Active {
		return nil, ErrInvalidUser
	}

	token, err := ts.NewAccessToken(user, strings.Split(claims.Scopes, " "), app, false, nil)
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

// NewInviteToken creates new invite token.
func (ts *JWTokenService) NewInviteToken(email, role string) (model.Token, error) {
	payload := make(map[string]interface{})
	// add payload data here
	if email != "" {
		payload["email"] = email
	}
	if role != "" {
		payload["role"] = role
	}

	now := ijwt.TimeFunc().Unix()

	lifespan := InviteTokenLifespan

	claims := &model.Claims{
		Payload: payload,
		Type:    model.TokenTypeInvite,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + lifespan,
			Issuer:    ts.issuer,
			// Subject:   u.ID(), //TODO: investigate why are we suppressing subject id from here?
			Audience: []string{"identifo"},
			IssuedAt: now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	token := model.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}
	return &model.JWToken{JWT: token, New: true}, nil
}

// NewResetToken creates new token for password resetting.
func (ts *JWTokenService) NewResetToken(userID string) (model.Token, error) {
	now := ijwt.TimeFunc().Unix()

	lifespan := ts.resetTokenLifespan

	claims := model.Claims{
		Type: model.TokenTypeReset,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   userID,
			Audience:  []string{"identifo"},
			IssuedAt:  now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	token := model.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &model.JWToken{JWT: token, New: true}, nil
}

// NewWebCookieToken creates new web cookie token.
func (ts *JWTokenService) NewWebCookieToken(u model.User) (model.Token, error) {
	if !u.Active {
		return nil, ErrInvalidUser
	}
	now := ijwt.TimeFunc().Unix()
	lifespan := ts.resetTokenLifespan

	claims := model.Claims{
		Type: model.TokenTypeWebCookie,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (now + lifespan),
			Issuer:    ts.issuer,
			Subject:   u.ID,
			Audience:  []string{"identifo"},
			IssuedAt:  now,
		},
	}

	var sm jwt.SigningMethod
	switch ts.algorithm {
	case model.TokenSignatureAlgorithmES256:
		sm = jwt.SigningMethodES256
	case model.TokenSignatureAlgorithmRS256:
		sm = jwt.SigningMethodRS256
	default:
		return nil, model.ErrWrongSignatureAlgorithm
	}

	token := model.NewTokenWithClaims(sm, ts.KeyID(), claims)
	if token == nil {
		return nil, ErrCreatingToken
	}

	return &model.JWToken{JWT: token, New: true}, nil
}

// String returns string representation of a token.
func (ts *JWTokenService) String(t model.Token) (string, error) {
	token, ok := t.(*model.JWToken)
	if !ok {
		return "", model.ErrTokenInvalid
	}

	if err := t.Validate(); err != nil {
		return "", err
	}
	if !token.New && !token.JWT.Valid {
		return "", model.ErrTokenInvalid
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
