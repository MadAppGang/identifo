package model

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type TokenType string

const (
	TokenTypeInvite     TokenType = "invite"     // TokenTypeInvite is an invite token type value.
	TokenTypeReset      TokenType = "reset"      // TokenTypeReset is an reset token type value.
	TokenTypeWebCookie  TokenType = "web-cookie" // TokenTypeWebCookie is a web-cookie token type value.
	TokenTypeAccess     TokenType = "access"     // TokenTypeAccess is an access token type.
	TokenTypeRefresh    TokenType = "refresh"    // TokenTypeRefresh is a refresh token type.
	TokenTypeManagement TokenType = "management" // TokenTypeManagement is a management token type for admin panel."
	TokenTypeID         TokenType = "id_token"   // id token type regarding oidc specification
	TokenTypeSignin     TokenType = "signin"     // signin token issues for user to sign in, etc to exchange for auth tokens. For example from admin panel I can send a link to user to siging to email with magic link.
	TokenTypeActor      TokenType = "actor"      // actor token is token impersonation. Admin could impersonated to be some of the users.

	// TODO: Deprecate it?
	// ! Deprecated: don't use it.
	TokenTypeTFAPreauth TokenType = "2fa-preauth" // TokenTypeTFAPreauth is an 2fa preauth token type.
	// TODO: Add other tokens, like admin, one, 2fa etc
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
	Payload() map[string]interface{}
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

func (t *JWToken) Claims() *Claims {
	return t.JWT.Claims.(*Claims)
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
// https://github.com/form3tech-oss/jwt-go/blob/master/cmd/jwt/app.go
