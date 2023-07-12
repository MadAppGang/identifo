package model

// TODO: implement key rotation
// TokenService is an abstract token manager.
type TokenService interface {
	// new methods
	NewToken(tokenType TokenType, u User, aud []string, fields []string, payload map[string]any) (Token, error)
	SignToken(token Token) (string, error)

	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string

	// keys management
	// replace the old private key with a new one
	SetPrivateKey(key any)

	// not using crypto.PublicKey here to avoid dependencies
	PublicKey() any
	KeyID() string
}
