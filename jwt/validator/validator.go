package validator

import (
	"errors"
	"os"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
)

var (
	// ErrTokenValidationNoExpiration is when the token does not have an expiration date.
	ErrTokenValidationNoExpiration = errors.New("Token is invalid, no expire date")
	// ErrTokenValidationExpired is when the token expiration date has passed
	ErrTokenValidationExpired = errors.New("Token is invalid, token has expired")
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
	// ErrorConfigurationMissingPublicKey is when public key is missing
	ErrorConfigurationMissingPublicKey = errors.New("Missing public key to decode the token from string")
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
	Validate(model.Token) error
	ValidateString(string) (model.Token, error)
}

// Config is a struct to set all the required params for Validator
type Config struct {
	Audience  []string
	Issuer    []string
	UserID    []string
	TokenType []string
	PublicKey interface{}
	// PubKeyEnvName environment variable for public key, could be empty if you want to use file insted
	PubKeyEnvName string
	// PubKeyFileName file path with public key, could be empty if you want to use env variable.
	PubKeyFileName string
	// PubKeyURL URL for well-known JWKS
	PubKeyURL string
	// should we always check audience for the token. If yes and audience is empty the validation will fail.
	IsAudienceRequired bool
	// should we always check iss for the token. If yes and iss is empty the validation will fail.
	IsIssuerRequired bool
}

// NewConfig creates and returns default config
func NewConfig() Config {
	return Config{
		TokenType:          []string{string(model.TokenTypeAccess)},
		IsAudienceRequired: true,
		IsIssuerRequired:   true,
	}
}

// NewValidator creates new JWT tokens validator.
// Arguments:
// - appID - application ID which have made the request, should be in audience field of JWT token.
// - issuer - this server name, should be the same as issuer of JWT token.
// - userID - user who have made the request. If this field is empty, we do not validate it.
func NewValidator(audience, issuer, userID, tokenType []string) Validator {
	return &validator{
		audience:   audience,
		issuer:     issuer,
		userID:     userID,
		tokenType:  tokenType,
		strictAud:  true,
		strictIss:  true,
		strictUser: false,
	}
}

// NewValidatorWithConfig creates new JWT tokens validator with public key from config file.
// Arguments:
// - appID - application ID which have made the request, should be in audience field of JWT token.
// - issuer - this server name, should be the same as issuer of JWT token.
// - userID - user who have made the request. If this field is empty, we do not validate it.
// - config - public key to parse the token.
func NewValidatorWithConfig(c Config) (Validator, error) {
	var key interface{}
	var err error = nil
	if len(c.PubKeyEnvName) > 0 {
		pk := os.Getenv(c.PubKeyEnvName)
		key, _, err = jwt.LoadPublicKeyFromString(pk)
	} else if len(c.PubKeyFileName) > 0 {
		key, _, err = jwt.LoadPublicKeyFromPEM(c.PubKeyFileName)
	}

	return &validator{
		audience:  c.Audience,
		issuer:    c.Issuer,
		userID:    c.UserID,
		tokenType: c.TokenType,
		strictAud: c.IsAudienceRequired,
		strictIss: c.IsIssuerRequired,
		publicKey: key,
	}, err
}

// TODO: implement initializer with JWKS URL .well-known

// validator is a JWT token validator.
type validator struct {
	audience   []string
	issuer     []string
	userID     []string
	tokenType  []string
	publicKey  interface{}
	strictIss  bool
	strictAud  bool
	strictUser bool
}

// Validate validates token.
func (v *validator) Validate(t model.Token) error {
	if t == nil {
		return model.ErrEmptyToken
	}
	// We assume the signature and standart claims were validated on parse.
	if err := t.Validate(); err != nil {
		return err
	}

	// We have already validated time based claims "exp, iat, nbf".
	// But, if any of the above claims are not in the token, it will still be considered a valid claim.
	// That's why these two fields are required: "exp, iat".
	token, ok := t.(*model.JWToken)
	if !ok {
		return model.ErrTokenInvalid
	}

	// Ensure the signature algorithm attack is not passing through.
	if token.JWT.Method.Alg() != SignatureAlgES && token.JWT.Method.Alg() != SignatureAlgRS {
		return model.ErrTokenInvalid
	}

	claims, ok := token.JWT.Claims.(*model.Claims)
	if !ok {
		return model.ErrTokenInvalid
	}

	if claims.ExpiresAt == 0 {
		return ErrTokenValidationNoExpiration
	}

	now := jwt.TimeFunc().Unix()
	if !claims.VerifyExpiresAt(now, true) {
		return ErrTokenValidationExpired
	}

	if !claims.VerifyIssuedAt(now, true) {
		return ErrTokenValidationNoIAT
	}

	// Validate Issuers
	if len(v.issuer) > 0 {
		valid := false
		for _, i := range v.issuer {
			if claims.VerifyIssuer(i, v.strictIss) {
				valid = true // at least one issues is valid, token is valid
				break
			}
		}
		if !valid {
			return ErrTokenValidationInvalidIssuer
		}
	}

	// Validate Audience
	if len(v.audience) > 0 {
		valid := false
		for _, i := range v.audience {
			if claims.VerifyAudience(i, v.strictAud) {
				valid = true // at least one audience is valid, token is valid
				break
			}
		}
		if !valid {
			return ErrTokenValidationInvalidAudience
		}
	}

	// Validate Users
	if len(v.userID) > 0 {
		valid := false
		for _, i := range v.userID {
			if (len(i) > 0) && (claims.Subject == i) {
				valid = true
			}
		}
		if !valid {
			return ErrTokenValidationInvalidSubject
		}
	}

	// Validate token type
	if len(v.tokenType) > 0 {
		valid := false
		for _, i := range v.tokenType {
			if token.Type() == i {
				valid = true
				break
			}
		}
		if !valid {
			return ErrorTokenValidationTokenTypeMismatch
		}
	}

	return nil
}

// ValidateString validates string representation of the token.
func (v *validator) ValidateString(t string) (model.Token, error) {
	if v.publicKey == nil {
		return nil, ErrorConfigurationMissingPublicKey
	}
	token, err := jwt.ParseTokenWithPublicKey(t, v.publicKey)
	if err != nil {
		return nil, err
	}
	return token, v.Validate(token)
}
