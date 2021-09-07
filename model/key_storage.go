package model

// Key names.
const (
	PublicKeyName  = "public.pem"
	PrivateKeyName = "private.pem"
)

type KeysPEM struct {
	Public  string `json:"public,omitempty"`
	Private string `json:"private,omitempty"`
}

// KeyStorage stores keys used for signing and verifying JWT tokens.
type KeyStorage interface {
	ReplaceKeys(keys JWTKeys) error
	LoadKeys(alg TokenSignatureAlgorithm) (JWTKeys, error)
	GetKeys() (KeysPEM, error)
}
