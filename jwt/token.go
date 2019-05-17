package jwt

import jwt "github.com/dgrijalva/jwt-go"

// Token is an abstract application token.
type Token interface {
	Validate() error
	UserID() string
	Type() string
	Payload() map[string]string
}

// NewTokenWithClaims generates new JWT token with claims and keyID.
func NewTokenWithClaims(method jwt.SigningMethod, kid string, claims jwt.Claims) *jwt.Token {
	return &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			"kid": kid,
		},
		Claims: claims,
		Method: method,
	}
}

// JWToken represents JWT token.
type JWToken struct {
	JWT *jwt.Token
	New bool
}

// Validate validates token data. Returns nil if all data is valid.
func (t *JWToken) Validate() error {
	if t.JWT == nil {
		return ErrEmptyToken
	}
	if !t.New && !t.JWT.Valid {
		return ErrTokenInvalid
	}
	return nil
}

// UserID returns user ID.
func (t *JWToken) UserID() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Payload returns token payload.
func (t *JWToken) Payload() map[string]string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return make(map[string]string)
	}
	return claims.Payload
}

// Type returns token type.
func (t *JWToken) Type() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Type
}

// Claims is an extended claims structure.
type Claims struct {
	Payload map[string]string `json:"payload,omitempty"`
	Scopes  string            `json:"scopes,omitempty"`
	Type    string            `json:"type,omitempty"`
	KeyID   string            `json:"kid,omitempty"` // optional keyID
	jwt.StandardClaims
}

// Full example of how to use JWT tokens:
// https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
