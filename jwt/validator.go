package jwt

import (
	"errors"
	"time"

	"github.com/madappgang/identifo/model"
)

var (
	ErrTokenValidationNoExpire        = errors.New("Token is invalid, no expire date")
	ErrTokenValidationNoIAT           = errors.New("Token is invalid, no issued at date")
	ErrTokenValidationInvalidIssuer   = errors.New("Token is invalid, issuer is invalid")
	ErrTokenValidationInvalidAudience = errors.New("Token is invalid, audience is invalid")
	ErrTokenValidationInvalidSubject  = errors.New("Token is invalid, subject is invalid")
)

const (
	//SignatureAlg is hardcoded signature algorithm
	//there is a number of options, we are stick to this value
	//see https://tools.ietf.org/html/rfc7516 for details
	SignatureAlg = "ES256"
)

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

//NewValidator creates new JWT validator
//appID - application ID who have made the request, should be in audience field of JWT token
//issues - this server name, should be the same as iss of JWT token
//userID - user, who have made the request, if the field is empty, we are not validating it
func NewValidator(appID, issuer, userID string) model.Validator {
	return &Validator{
		appID:  appID,
		issuer: issuer,
		userID: userID,
	}
}

//Validator JWT token validator
type Validator struct {
	appID  string
	issuer string
	userID string
}

//Validate validates token
func (v *Validator) Validate(t model.Token) error {
	if t == nil {
		return ErrEmptyToken
	}
	//we assume the signature and standart claims were validated on parse
	if err := t.Validate(); err != nil {
		return err
	}
	// We have already have validated time based claims "exp, iat, nbf".
	// But, if any of the above claims are not in the token, it will still
	// be considered a valid claim.
	// That's why all these two fields are required: "exp, iat"
	token, ok := t.(*Token)
	if !ok {
		return ErrTokenInvalid
	}

	//check the signature algorithm attack is not passing through
	if token.JWT.Method.Alg() != SignatureAlg {
		return ErrTokenInvalid
	}

	claims, ok := token.JWT.Claims.(*Claims)
	if !ok {
		return ErrTokenInvalid
	}

	now := TimeFunc().Unix()
	if !claims.VerifyExpiresAt(now, true) {
		return ErrTokenValidationNoExpire
	}

	if !claims.VerifyIssuedAt(now, true) {
		return ErrTokenValidationNoIAT
	}

	if !claims.VerifyAudience(v.appID, true) {
		return ErrTokenValidationInvalidAudience
	}

	if !claims.VerifyIssuer(v.issuer, true) {
		return ErrTokenValidationInvalidIssuer
	}

	if (len(v.userID) > 0) && (claims.Subject != v.userID) {
		return ErrTokenValidationInvalidSubject
	}

	return nil
}
