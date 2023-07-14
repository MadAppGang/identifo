package model

import (
	"encoding/json"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/tools/xmaps"
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

func (t TokenType) String() string {
	return string(t)
}

// TokenWithClaims generates new JWT token with claims and keyID.
func TokenWithClaims(method jwt.SigningMethod, kid string, claims jwt.Claims) *JWToken {
	return &JWToken{
		Token: jwt.Token{
			Header: map[string]any{
				"typ": "JWT",
				"alg": method.Alg(),
				"kid": kid,
			},
			Claims: claims,
			Method: method,
		},
		New: true,
	}
}

// JWToken represents JWT token.
type JWToken struct {
	jwt.Token
	New bool
}

func (t *JWToken) FullClaims() Claims {
	claims, _ := t.Claims.(*Claims)
	return *claims
}

// Validate validates token data. Returns nil if all data is valid.
func (t *JWToken) Validate() error {
	// pass jwt lib token validation as it has not been parsed, it was constructed.
	if t.New {
		return nil
	}

	// the token is invalid by jwt validator while parsing
	if !t.Token.Valid {
		return l.ErrorTokenInvalid
	}

	return nil
}

// UserID returns user ID.
func (t *JWToken) UserID() string {
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Payload returns token payload.
func (t *JWToken) Payload() map[string]any {
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil
	}
	return claims.Payload
}

// Type returns token type.
func (t *JWToken) Type() TokenType {
	claims, ok := t.Claims.(Claims)
	if !ok {
		return ""
	}
	return TokenType(claims.Type)
}

// ExpiresAt standard token claim
func (t *JWToken) ExpiresAt() time.Time {
	claims, ok := t.Claims.(*Claims)
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
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.ID
}

// IssuedAt standard token claim
func (t *JWToken) IssuedAt() time.Time {
	claims, ok := t.Claims.(*Claims)
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
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Issuer
}

// NotBefore standard token claim
func (t *JWToken) NotBefore() time.Time {
	claims, ok := t.Claims.(*Claims)
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
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

// Claims is an extended claims structure.
type Claims struct {
	Payload map[string]any `json:"payload,omitempty"`
	Type    string         `json:"type,omitempty"`
	KeyID   string         `json:"kid,omitempty"` // optional keyID
	jwt.RegisteredClaims
}

// MarshalJSON marshal everything, flattering to the root level.
func (c Claims) MarshalJSON() ([]byte, error) {
	m := xmaps.LowercaseKeys(c.Payload)
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
	if len(c.Type) > 0 {
		m["type"] = c.Type
	}
	if len(c.KeyID) > 0 {
		m["kid"] = c.KeyID
	}
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
	exclude := map[string]bool{"iss": true, "sub": true, "aud": true, "exp": true, "nbf": true, "iat": true, "jti": true, "kid": true, "type": true}
	c.Payload = map[string]any{}
	for k, v := range pc {
		lk := strings.ToLower(k)
		if !exclude[lk] {
			// try to convert slice of any to slice of strings for roles
			if strings.HasPrefix(lk, RoleScopePrefix) {
				c.Payload[lk] = toSliceString(v)
			} else {
				c.Payload[lk] = v
			}
		}
	}
	if pc["kid"] != nil {
		c.KeyID = pc["kid"].(string)
	}
	if pc["type"] != nil {
		c.Type = pc["type"].(string)
	}

	return nil
}

// if it fail on any conversion - return the original value
func toSliceString(s any) any {
	sl, ok := s.([]any)
	if !ok {
		return s
	}
	r := []string{}
	for _, v := range sl {
		vs, ok := v.(string)
		if !ok {
			return s
		}
		r = append(r, vs)
	}
	return r
}

// Full example of how to use JWT tokens:
// https://github.com/form3tech-oss/jwt-go/blob/master/cmd/jwt/app.go
