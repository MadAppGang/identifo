package model

import (
	"encoding/json"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/maps"
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
)

// StandardTokenClaims structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
type StandardTokenClaims interface {
	Audience() []string
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
		Header: map[string]any{
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
	jwt.Token
	New bool
}

// Validate validates token data. Returns nil if all data is valid.
func (t *JWToken) Validate() error {
	if !t.New {
		return t.Validate()
	}
	return nil
}

// UserID returns user ID.
func (t *JWToken) UserID() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Payload returns token payload.
func (t *JWToken) Payload() map[string]interface{} {
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil
	}
	return claims.Payload
}

// Type returns token type.
func (t *JWToken) Type() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.Type
}

// Audience standard token claim
func (t *JWToken) Audience() []string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return nil
	}
	return claims.Audience
}

// ExpiresAt standard token claim
func (t *JWToken) ExpiresAt() time.Time {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return time.Time{}
	}
	if claims.ExpiresAt == nil {
		return time.Time{}
	}
	return (*claims.ExpiresAt).Time
}

// ID standard token claim
func (t *JWToken) ID() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.ID
}

// func (t *JWToken) Claims() *Claims {
// 	return t.JWT.Claims.(*Claims)
// }

// IssuedAt standard token claim
func (t *JWToken) IssuedAt() time.Time {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return time.Time{}
	}
	if claims.IssuedAt == nil {
		return time.Time{}
	}
	return (*claims.IssuedAt).Time
}

// Issuer standard token claim
func (t *JWToken) Issuer() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.Issuer
}

// NotBefore standard token claim
func (t *JWToken) NotBefore() time.Time {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return time.Time{}
	}
	if claims.NotBefore == nil {
		return time.Time{}
	}
	return (*claims.NotBefore).Time
}

// Subject standard token claim
func (t *JWToken) Subject() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Scopes standard token claim
func (t *JWToken) Scopes() string {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return claims.Scopes
}

// Claims is an extended claims structure.
type Claims struct {
	Payload map[string]any `json:"payload,omitempty"`
	Scopes  string         `json:"scopes,omitempty"`
	Type    string         `json:"type,omitempty"`
	KeyID   string         `json:"kid,omitempty"` // optional keyID
	jwt.RegisteredClaims
}

// MarshalJSON marshal everything, flattering to the root level.
func (c *Claims) MarshalJSON() ([]byte, error) {
	m := maps.Clone(c.Payload)
	type proxy jwt.RegisteredClaims
	var p proxy = proxy(c.RegisteredClaims)
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	var rcm map[string]any
	err = json.Unmarshal(b, &rcm)
	if err != nil {
		return nil, err
	}
	maps.Copy(m, rcm)
	return json.Marshal(&m)
}

// UnmarshalJSON unmarshals JWT token from flat structure to it's original structure.
func (c *Claims) UnmarshalJSON(data []byte) error {
	var rc jwt.RegisteredClaims
	if err := json.Unmarshal(data, &rc); err != nil {
		return err
	}
	c.RegisteredClaims = rc

	var pc map[string]any
	if err := json.Unmarshal(data, &pc); err != nil {
		return err
	}
	exclude := map[string]bool{"iss": true, "sub": true, "aud": true, "exp": true, "nbf": true, "iat": true, "jti": true}
	c.Payload = map[string]any{}
	for k, v := range pc {
		if !exclude[k] {
			c.Payload[k] = v
		}
	}

	return nil
}

// Full example of how to use JWT tokens:
// https://github.com/form3tech-oss/jwt-go/blob/master/cmd/jwt/app.go
