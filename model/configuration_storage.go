package model

// ConfigurationStorage stores server configuration.
type ConfigurationStorage interface {
	WriteConfig(ServerSettings) error
	LoadServerSettings(forceReload bool) (ServerSettings, error)
	InsertKeys(keys *JWTKeys) error
	LoadKeys(TokenSignatureAlgorithm) (*JWTKeys, error)
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
	LoadKeys(alg TokenSignatureAlgorithm) (*JWTKeys, error)
}
