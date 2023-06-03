package model

const (
	// OfflineScope is a scope value to request refresh token.
	OfflineScope = "offline"
)

// TokenService is an abstract token manager.
type TokenService interface {
	NewAccessToken(u User, scopes []string, app AppData, requireTFA bool, tokenPayload map[string]interface{}) (Token, error)
	NewRefreshToken(u User, scopes []string, app AppData) (Token, error)
	RefreshAccessToken(token Token, tokenPayload map[string]interface{}) (Token, error)
	NewInviteToken(email, role, audience string, data map[string]interface{}) (Token, error)
	NewResetToken(userID string) (Token, error)
	NewWebCookieToken(u User) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string
	WebCookieTokenLifespan() int64

	// keys management
	// replace the old private key with a new one
	SetPrivateKey(key interface{})
	PrivateKey() interface{}
	// not using crypto.PublicKey here to avoid dependencies
	PublicKey() interface{}
	KeyID() string
}
