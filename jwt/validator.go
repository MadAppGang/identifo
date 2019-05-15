package jwt

import (
	"errors"
	"time"
)

var (
	// ErrTokenValidationNoExpiration is when the token does not have an expiration date.
	ErrTokenValidationNoExpiration = errors.New("Token is invalid, no expire date")
	// ErrTokenValidationNoIAT is when IAT verification fails.
	ErrTokenValidationNoIAT = errors.New("Token is invalid, no issued at date")
	// ErrTokenValidationInvalidIssuer is when the token has invalid issuer.
	ErrTokenValidationInvalidIssuer = errors.New("Token is invalid, issuer is invalid")
	// ErrTokenValidationInvalidAudience is when the token has invalid audience.
	ErrTokenValidationInvalidAudience = errors.New("Token is invalid, audience is invalid")
	// ErrTokenValidationInvalidSubject is when subject claim is invalid.
	ErrTokenValidationInvalidSubject = errors.New("Token is invalid, subject is invalid")
	// ErrorTokenValidationTokenTypeMismatch is when the token has invalid type.
	ErrorTokenValidationTokenTypeMismatch = errors.New("Token is invalid, type is invalid")
)

const (
	// SignatureAlgES is a hardcoded ES256 signature algorithm.
	// There is a number of options, we are stick to this value.
	// See https://tools.ietf.org/html/rfc7516 for details.
	SignatureAlgES = "ES256"
	// SignatureAlgRS is a hardcoded RS256 signature algorithm.
	SignatureAlgRS = "RS256"
)

// Validator is an abstract token validator.
type Validator interface {
	Validate(Token) error
}

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value. This is useful for testing or if your
// server uses a time zone different from your tokens'.
var TimeFunc = time.Now

// NewDefaultValidator creates new JWT tokens validator.
// Arguments:
// - appID - application ID which have made the request, should be in audience field of JWT token.
// - issuer - this server name, should be the same as issuer of JWT token.
// - userID - user who have made the request. If this field is empty, we do not validate it.
func NewDefaultValidator(audience, issuer, userID, tokenType string) Validator {
	return &DefaultValidator{
		audience:  audience,
		issuer:    issuer,
		userID:    userID,
		tokenType: tokenType,
	}
}

// DefaultValidator is a default JWT token validator.
type DefaultValidator struct {
	audience  string
	issuer    string
	userID    string
	tokenType string
}

// Validate validates token.
func (v *DefaultValidator) Validate(t Token) error {
	if t == nil {
		return ErrEmptyToken
	}
	// We assume the signature and standart claims were validated on parse.
	if err := t.Validate(); err != nil {
		return err
	}

	// We have already validated time based claims "exp, iat, nbf".
	// But, if any of the above claims are not in the token, it will still be considered a valid claim.
	// That's why these two fields are required: "exp, iat".
	token, ok := t.(*DefaultToken)
	if !ok {
		return ErrTokenInvalid
	}

	// Ensure the signature algorithm attack is not passing through.
	if token.JWT.Method.Alg() != SignatureAlgES && token.JWT.Method.Alg() != SignatureAlgRS {
		return ErrTokenInvalid
	}

	claims, ok := token.JWT.Claims.(*Claims)
	if !ok {
		return ErrTokenInvalid
	}

	now := TimeFunc().Unix()
	if !claims.VerifyExpiresAt(now, true) {
		return ErrTokenValidationNoExpiration
	}

	if !claims.VerifyIssuedAt(now, true) {
		return ErrTokenValidationNoIAT
	}

	if !claims.VerifyAudience(v.audience, true) {
		return ErrTokenValidationInvalidAudience
	}

	if !claims.VerifyIssuer(v.issuer, true) {
		return ErrTokenValidationInvalidIssuer
	}

	if (len(v.userID) > 0) && (claims.Subject != v.userID) {
		return ErrTokenValidationInvalidSubject
	}

	if token.Type() != v.tokenType {
		return ErrorTokenValidationTokenTypeMismatch
	}

	return nil
}

// ValidateString validates string representation of the token.
func (v *DefaultValidator) ValidateString(t string, publicKey interface{}) error {
	token, err := ParseTokenWithPublicKey(t, publicKey)
	if err != nil {
		return err
	}
	return v.Validate(token)
}
