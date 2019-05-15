package tokensrvc

import (
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

const (
	// OfflineScope is a scope value to request refresh token.
	OfflineScope = "offline"
	// RefrestTokenType is a refresh token type value.
	RefrestTokenType = "refresh"
	// InviteTokenType is an invite token type value.
	InviteTokenType = "invite"
	// AccessTokenType is an access token type value.
	AccessTokenType = "access"
	// ResetTokenType is a reset password token type value.
	ResetTokenType = "reset"
	// WebCookieTokenType is a web-cookie token type value.
	WebCookieTokenType = "web-cookie"
)

// TokenService is an abstract token manager.
type TokenService interface {
	NewToken(u model.User, scopes []string, app model.AppData) (jwt.Token, error)
	NewRefreshToken(u model.User, scopes []string, app model.AppData) (jwt.Token, error)
	NewInviteToken() (jwt.Token, error)
	NewResetToken(userID string) (jwt.Token, error)
	RefreshToken(token jwt.Token) (jwt.Token, error)
	NewWebCookieToken(u model.User) (jwt.Token, error)
	Parse(string) (jwt.Token, error)
	String(jwt.Token) (string, error)
	Issuer() string
	Algorithm() string
	WebCookieTokenLifespan() int64
	PublicKey() interface{} // we are not using crypto.PublicKey here to avoid dependencies
	KeyID() string
}
