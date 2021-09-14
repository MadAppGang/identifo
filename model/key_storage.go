package model

// Key names.
const PrivateKeyName = "private.pem"

// KeyStorage stores keys used for signing and verifying JWT tokens.
type KeyStorage interface {
	ReplaceKey(keyPEM []byte) error
	LoadPrivateKey() (interface{}, error)
}
