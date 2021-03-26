package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// StandardTokenClaims structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
type StandardTokenClaims interface {
	Audience() string
	ExpiresAt() time.Time
	ID() string
	IssuedAt() time.Time
	Issuer() string
	NotBefore() time.Time
	Subject() string
}

// Token is an abstract application token.
type Token interface {
	StandardTokenClaims
	Validate() error
	UserID() string
	Type() string
	Scopes() string
    Payload() map[string]intereface{}
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
func (t *JWToken) Payload() map[string]interface{} {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return make(map[string]interface{})
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

// Audience standard token claim
func (t *JWToken) Audience() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Audience
}

// ExpiresAt standard token claim
func (t *JWToken) ExpiresAt() time.Time {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return time.Time{}
	}
	return time.Unix(claims.ExpiresAt, 0)
}

// ID standard token claim
func (t *JWToken) ID() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Id
}

// IssuedAt standard token claim
func (t *JWToken) IssuedAt() time.Time {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return time.Time{}
	}
	return time.Unix(claims.IssuedAt, 0)
}

// Issuer standard token claim
func (t *JWToken) Issuer() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Issuer
}

// NotBefore standard token claim
func (t *JWToken) NotBefore() time.Time {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return time.Time{}
	}
	return time.Unix(claims.NotBefore, 0)
}

// Subject standard token claim
func (t *JWToken) Subject() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Scopes standard token claim
func (t *JWToken) Scopes() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Scopes
}

// Claims is an extended claims structure.
type Claims struct {
	Payload map[string]interface{} `json:"payload,omitempty"`
	Scopes  string                 `json:"scopes,omitempty"`
	Type    string                 `json:"type,omitempty"`
	KeyID   string                 `json:"kid,omitempty"` // optional keyID
	jwt.StandardClaims
}

// Full example of how to use JWT tokens:
// https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
