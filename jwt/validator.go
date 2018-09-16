package jwt

import (
	"errors"
	"time"

	"github.com/madappgang/identifo/model"
)

var (
	ErrTokenValidationNoExpire = errors.New("Token is invalid, no expire date")
	ErrTokenValidationNoIAT    = errors.New("Token is invalid, no issued at date")
)

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

//NewValidator creates new JWT validator
func NewValidator(appID, issuer string) model.Validator {
	v := Validator{}
	v.appID = appID
	v.issuer = issuer
	return &v
}

//Validator JWT token validator
type Validator struct {
	appID  string
	issuer string
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
	claims, ok := token.JWT.Claims.(*Claims)
	if !ok {
		return ErrTokenInvalid
	}

	now := TimeFunc().Unix()
	if claims.VerifyExpiresAt(now, true) == false {
		return ErrTokenValidationNoExpire
	}

	if claims.VerifyIssuedAt(now, true) == false {
		return ErrTokenValidationNoIAT
	}

	///Audience //servers who will use it token
	///Issuer   //who issued, should be this server
	///Subject //user ID of the token
	///

	return nil
}
