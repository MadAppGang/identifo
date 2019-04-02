package model

import (
	"encoding/json"
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
	TokenServiceAlgorithmES256 TokenServiceAlgorithm = iota + 1
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

// String implements Stringer.
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

// MarshalJSON implements json.Marshaller.
func (alg TokenServiceAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(alg.String())
}

// UnmarshalJSON implements json.Unmarshaller.
func (alg *TokenServiceAlgorithm) UnmarshalJSON(data []byte) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	algorithm, ok := map[string]TokenServiceAlgorithm{
		"es256": TokenServiceAlgorithmES256,
		"rs256": TokenServiceAlgorithmRS256,
		"auto":  TokenServiceAlgorithmAuto}[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenServiceAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}

// MarshalYAML implements yaml.Marshaller.
func (alg TokenServiceAlgorithm) MarshalYAML() (interface{}, error) {
	return alg.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (alg *TokenServiceAlgorithm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := unmarshal(&aux); err != nil {
		return err
	}

	algorithm, ok := map[string]TokenServiceAlgorithm{
		"es256": TokenServiceAlgorithmES256,
		"rs256": TokenServiceAlgorithmRS256,
		"auto":  TokenServiceAlgorithmAuto}[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenServiceAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}
