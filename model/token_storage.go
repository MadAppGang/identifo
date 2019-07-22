package model

// TokenStorage is a storage for issued refresh tokens.
type TokenStorage interface {
	SaveToken(token string) error
	HasToken(token string) bool
	RevokeToken(token string) error
	Close()
}

// TokenBlacklist is a storage for blacklisted tokens.
type TokenBlacklist interface {
	IsBlacklisted(token string) bool
	Add(token string) error
	Close()
}
