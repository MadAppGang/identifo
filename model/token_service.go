package model

import (
	"fmt"
)

const (
	// OfflineScope is a scope value to request refresh token.
	OfflineScope = "offline"
	// RefrestTokenType is a refresh token type value.
	RefrestTokenType = "refresh"
	// AccessTokenType is an access token type value.
	AccessTokenType = "access"
	// ResetTokenType is a reset password token type value.
	ResetTokenType = "reset"
	// WebCookieTokenType is a web-cookie token type value.
	WebCookieTokenType = "web-cookie"
)

// TokenServiceAlgorithm - we support only two now.
type TokenServiceAlgorithm int

const (
	// TokenServiceAlgorithmES256 is a ES256 signature.
	TokenServiceAlgorithmES256 TokenServiceAlgorithm = iota
	// TokenServiceAlgorithmRS256 is a RS256 signature.
	TokenServiceAlgorithmRS256
	// TokenServiceAlgorithmAuto tries to detect algorithm on the fly.
	TokenServiceAlgorithmAuto
)

// TokenService manages tokens abstraction layer.
type TokenService interface {
	// NewToken creates new access token for the user.
	NewToken(u User, scopes []string, app AppData) (Token, error)
	// NewRefreshToken creates new refresh token for the user.
	NewRefreshToken(u User, scopes []string, app AppData) (Token, error)
	// NewRestToken creates new reset password token.
	NewResetToken(userID string) (Token, error)
	// RefreshToken issues the new access token with access token.
	RefreshToken(token Token) (Token, error)
	NewWebCookieToken(u User) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string
	WebCookieTokenLifespan() int64
	PublicKey() interface{} // we are not using crypto.PublicKey here to avoid dependencies
	KeyID() string
}

// Token is an app token to give user chan
type Token interface {
	Validate() error
	UserID() string
	Type() string
	Payload() map[string]string
}

// Validator validates token with external requester.
type Validator interface {
	Validate(Token) error
}

// TokenMapping is a service for matching tokens to services.
type TokenMapping interface{}

func (alg TokenServiceAlgorithm) String() string {
	switch alg {
	case TokenServiceAlgorithmES256:
		return "es256"
	case TokenServiceAlgorithmRS256:
		return "rs256"
	case TokenServiceAlgorithmAuto:
		return "auto"
	default:
		return fmt.Sprintf("TokenServiceAlgorithm(%d)", alg)
	}
}
