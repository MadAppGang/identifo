package validator

import (
	"errors"
	"os"
	"time"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xslices"
	"golang.org/x/exp/slices"
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
	// PubKeyEnvName environment variable for public key, could be empty if you want to use file instead.
	PubKeyEnvName string
	// PubKeyFileName file path with public key, could be empty if you want to use env variable.
	PubKeyFileName string
	// PubKeyURL URL for well-known JWKS.
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
	subject    []string
	publicKey  interface{}
	strictIss  bool
	strictAud  bool
	strictUser bool
}

// Validate validates token.
func (v *validator) Validate(t model.Token) error {
	var errs error
	if t == nil {
		return l.ErrorTokenInvalid
	}
	// We assume the signature and standard claims were validated on parse.
	if err := t.Validate(); err != nil {
		return err
	}

	// We have already validated time based claims "exp, iat, nbf".
	// But, if any of the above claims are not in the token, it will still be considered a valid claim.
	// That's why these two fields are required: "exp, iat".
	token, ok := t.(*model.JWToken)
	if !ok {
		return l.ErrorTokenInvalid
	}

	// Ensure the signature algorithm attack is not passing through.
	if token.Method.Alg() != SignatureAlgES && token.Method.Alg() != SignatureAlgRS {
		errors.Join(errs, l.ErrorValidatingTokenMethod)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		errors.Join(errs, l.ErrorValidationTokenClaims)
	}

	now := jwt.TimeFunc()
	if err := v.verifyExpiresAt(*claims, now, true); err != nil {
		errors.Join(errs, err)
	}

	if err := v.verifyNotBefore(*claims, now, false); err != nil {
		errors.Join(errs, err)
	}

	if err := v.verifyIssuedAt(*claims, now, true); err != nil {
		errors.Join(errs, err)
	}

	// Validate Audience
	if err := v.verifyAudience(*claims, v.audience, v.strictAud); err != nil {
		errors.Join(errs, err)
	}
	// Validate Issuers
	if err := v.verifyIssuer(*claims, v.issuer, v.strictIss); err != nil {
		errors.Join(errs, err)
	}

	// Validate subject
	if err := v.verifySubject(*claims, v.userID, v.strictUser); err != nil {
		errors.Join(errs, err)
	}

	// Validate token type
	if err := v.verifyTokenType(*claims, v.tokenType, true); err != nil {
		errors.Join(errs, err)
	}

	return errs
}

// ValidateString validates string representation of the token.
func (v *validator) ValidateString(t string) (model.Token, error) {
	if v.publicKey == nil {
		return nil, l.ErrorServiceTokenValidatorNoPublicKey
	}
	token, err := jwt.ParseTokenWithPublicKey(t, v.publicKey)
	if err != nil {
		return nil, err
	}
	return token, v.Validate(token)
}

func (v *validator) verifyExpiresAt(claims model.Claims, cmp time.Time, required bool) error {
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}

	if exp == nil {
		return errorIfRequired(required, "exp")
	}

	return errorIfFalse(cmp.Before((exp.Time)), l.ErrorValidationTokenExpired)
}

func (v *validator) verifyNotBefore(claims model.Claims, cmp time.Time, required bool) error {
	nbf, err := claims.GetNotBefore()
	if err != nil {
		return err
	}

	if nbf == nil {
		return errorIfRequired(required, "nbf")
	}

	return errorIfFalse(!cmp.Before((nbf.Time)), l.ErrorValidationTokenNotValidYet)
}

func (v *validator) verifyIssuedAt(claims model.Claims, cmp time.Time, required bool) error {
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return err
	}

	if iat == nil {
		return errorIfRequired(required, "iat")
	}

	return errorIfFalse(!cmp.Before((iat.Time)), l.ErrorValidationTokenIssuedInFuture)
}

func (v *validator) verifyAudience(claims model.Claims, expected []string, required bool) error {
	aud, err := claims.GetAudience()
	if err != nil {
		return err
	}

	if len(aud) == 0 {
		return errorIfRequired(required, "aud")
	}

	found := xslices.Intersect(expected, aud)
	// nothing found
	if len(found) == 0 {
		return errorIfRequired(required, "aud")
	}

	return nil
}

func (v *validator) verifyIssuer(claims model.Claims, expected []string, required bool) error {
	iss, err := claims.GetIssuer()
	if err != nil {
		return err
	}

	if len(iss) == 0 {
		return errorIfRequired(required, "iss")
	}

	return errorIfFalse(slices.Contains(expected, iss), l.ErrorValidationTokenInvalidIssuer)
}

func (v *validator) verifySubject(claims model.Claims, expected []string, required bool) error {
	sub, err := claims.GetSubject()
	if err != nil {
		return err
	}

	if len(sub) == 0 {
		return errorIfRequired(required, "sub")
	}

	return errorIfFalse(slices.Contains(expected, sub), l.ErrorValidationTokenInvalidSubject)
}

func (v *validator) verifyTokenType(claims model.Claims, expected []string, required bool) error {
	if len(claims.Type) == 0 {
		return errorIfRequired(required, "sub")
	}

	return errorIfFalse(slices.Contains(expected, claims.Type), l.ErrorValidationTokenInvalidType)
}

// errorIfRequired returns an ErrorValidationTokenMissingClaim error if required is
// true. Otherwise, nil is returned.
func errorIfRequired(required bool, claim string) error {
	if required {
		return l.LocalizedError{
			ErrID:   l.ErrorValidationTokenMissingClaim,
			Details: []any{claim},
		}
	} else {
		return nil
	}
}

// errorIfFalse returns the error specified in err, if the value is true.
// Otherwise, nil is returned.
func errorIfFalse(value bool, err error) error {
	if value {
		return nil
	} else {
		return err
	}
}
