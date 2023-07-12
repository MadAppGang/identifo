package service

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	j "github.com/madappgang/identifo/v2/jwt"
	jv "github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xmaps"
	"golang.org/x/exp/maps"
)

// NewJWTokenService returns new JWT token service.
// Arguments:
// - privateKeyPath - the path to the private key in pem format. Please keep it in a secret place.
// - publicKeyPath - the path to the public key.
func NewJWTokenService(privateKey any, issuer string, settings model.SecurityServerSettings) (model.TokenService, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key is empty")
	}

	t := &JWTokenService{
		iss:      issuer,
		pk:       privateKey,
		settings: settings,
	}

	return t, nil
}

// JWTokenService is a JWT token service.
type JWTokenService struct {
	pk       any // *ecdsa.PrivateKey, or *rsa.PrivateKey
	settings model.SecurityServerSettings
	iss      string
	aCache   string // algorithm cache
	pkCache  any    // *ecdsa.PublicKey, or *rsa.PublicKey
}

// Issuer returns token issuer name.
func (ts *JWTokenService) Issuer() string {
	return ts.iss
}

// Issuer returns token issuer name.
func (ts *JWTokenService) PrivateKey() any {
	return ts.pk
}

func (ts *JWTokenService) NewToken(tokenType model.TokenType, u model.User, aud []string, fields []string, payload map[string]any) (model.JWToken, error) {
	// we have to collect all payloads to one map
	userPayload := xmaps.FieldsToMap(u)
	userPayload = xmaps.FilterMap(userPayload, fields)
	if payload == nil {
		payload = map[string]any{}
	}
	maps.Copy(payload, userPayload)
	lifespan := ts.settings.TokenLifetime(tokenType)
	ia := jwt.NewNumericDate(j.TimeFunc())
	exp := ia.Add(time.Minute * time.Duration(lifespan))

	claims := model.Claims{
		Payload: maps.Clone(payload),
		Type:    string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    ts.iss,
			Subject:   u.ID,
			Audience:  aud,
			IssuedAt:  ia,
		},
	}

	sm := ts.jwtMethod()
	if sm == nil {
		return model.JWToken{}, l.ErrorTokenMethodInvalid
	}

	token := model.TokenWithClaims(sm, ts.KeyID(), claims)
	return token, nil
}

func (ts *JWTokenService) SignToken(t model.JWToken) (string, error) {
	if err := t.Validate(); err != nil {
		return "", l.LocalizedError{
			ErrID:   l.ErrorValidatingToken,
			Details: []any{err},
		}
	}

	str, err := t.SignedString(ts.pk)
	if err != nil {
		return "", err
	}
	return str, nil
}

// Algorithm  returns signature algorithm.
func (ts *JWTokenService) Algorithm() string {
	if len(ts.aCache) > 0 {
		return ts.aCache
	}

	switch ts.pk.(type) {
	case *rsa.PrivateKey:
		ts.aCache = "RS256"
		return "RS256"
	case *ecdsa.PrivateKey:
		ts.aCache = "ES256"
		return "ES256"
	default:
		return ""
	}
}

func (ts *JWTokenService) jwtMethod() jwt.SigningMethod {
	switch ts.Algorithm() {
	case "ES256":
		return jwt.SigningMethodES256
	case "RS256":
		return jwt.SigningMethodRS256
	default:
		return nil
	}
}

// PublicKey returns public key.
func (ts *JWTokenService) PublicKey() any {
	if ts.pkCache != nil {
		return ts.pkCache
	}

	switch ts.pk.(type) {
	case *rsa.PrivateKey:
		pk := ts.pk.(*rsa.PrivateKey)
		ts.pkCache = pk.Public()
	case *ecdsa.PrivateKey:
		pk := ts.pk.(*ecdsa.PrivateKey)
		ts.pkCache = pk.Public()
	default:
		return nil
	}
	return ts.pkCache
}

func (ts *JWTokenService) SetPrivateKey(key any) {
	// Log event
	ts.pk = key
	ts.pkCache = nil
	ts.aCache = ""
}

// KeyID returns public key ID, using SHA-1 fingerprint.
func (ts *JWTokenService) KeyID() string {
	pk := ts.PublicKey()
	if pk != nil {
		if der, err := x509.MarshalPKIXPublicKey(pk); err == nil {
			s := sha1.Sum(der)
			return base64.RawURLEncoding.EncodeToString(s[:]) // slice from [20]byte
		}
	}
	return ""
}

// Parse parses token data from the string representation.
func (ts *JWTokenService) Parse(s string) (model.JWToken, error) {
	tokenString := strings.TrimSpace(s)
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (any, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counterpart to verify them.
		// TODO: add multi-key supports
		return ts.PublicKey(), nil
	})
	if err != nil {
		return model.JWToken{}, err
	}

	return model.JWToken{Token: *token}, nil
}

// ValidateTokenString parses token and validates it.
func (ts *JWTokenService) ValidateTokenString(tstr string, v jv.Validator, tokenType string) (model.JWToken, error) {
	token, err := ts.Parse(tstr)
	if err != nil {
		return model.JWToken{}, err
	}

	if err := v.Validate(token); err != nil {
		return model.JWToken{}, err
	}

	if token.Type() != tokenType {
		return model.JWToken{}, err
	}

	return token, nil
}
