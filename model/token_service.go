package model

// TODO: refactor to reduce number of methods
// TODO: implement key rotation
// TokenService is an abstract token manager.
type TokenService interface {
	// new methods
	NewToken(tokenType TokenType, u User, fields []string, payload map[string]any) (Token, error)
	SignToken(token Token) (string, error)

	// // old methods
	// NewAccessToken(u User, scopes []string, app AppData, requireTFA bool, tokenPayload map[string]interface{}) (Token, error)
	// NewRefreshToken(u User, scopes []string, app AppData) (Token, error)
	// RefreshAccessToken(token Token) (Token, error)
	// NewInviteToken(email, role, audience string, data map[string]interface{}) (Token, error)
	// NewResetToken(userID string) (Token, error)
	// // NewToken(tokenType TokenType, userID string, payload []any) (Token, error)
	// NewWebCookieToken(u User) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string
	WebCookieTokenLifespan() int64

	// keys management
	// replace the old private key with a new one
	SetPrivateKey(key any)

	// not using crypto.PublicKey here to avoid dependencies
	PublicKey() any
	KeyID() string
}
