package model

import (
	ijwt "github.com/madappgang/identifo/jwt"
)

// ConfigurationStorage stores server configuration.
type ConfigurationStorage interface {
	InsertConfig(key string, value interface{}) error
	LoadServerSettings(*ServerSettings) error
	InsertKeys(keys *JWTKeys) error
	LoadKeys(ijwt.TokenSignatureAlgorithm) (*JWTKeys, error)
	GetUpdateChan() chan interface{}
	CloseUpdateChan()
}

// Key names.
const (
	PublicKeyName  = "public.pem"
	PrivateKeyName = "private.pem"
)

// KeyStorage stores keys used for signing and verifying JWT tokens.
type KeyStorage interface {
	InsertKeys(keys *JWTKeys) error
	LoadKeys(alg ijwt.TokenSignatureAlgorithm) (*JWTKeys, error)
}
