package jwt

import (
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
	NewToken(u model.User, scopes []string, app model.AppData) (Token, error)
	NewRefreshToken(u model.User, scopes []string, app model.AppData) (Token, error)
	NewInviteToken() (Token, error)
	NewResetToken(userID string) (Token, error)
	RefreshToken(token Token) (Token, error)
	NewWebCookieToken(u model.User) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string
	WebCookieTokenLifespan() int64
	PublicKey() interface{} // we are not using crypto.PublicKey here to avoid dependencies
	KeyID() string
}
