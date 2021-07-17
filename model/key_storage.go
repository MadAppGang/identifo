package model

// Key names.
const (
	PublicKeyName  = "public.pem"
	PrivateKeyName = "private.pem"
)

// KeyStorage stores keys used for signing and verifying JWT tokens.
type KeyStorage interface {
	InsertKeys(keys JWTKeys) error
	LoadKeys(alg TokenSignatureAlgorithm) (JWTKeys, error)
}
