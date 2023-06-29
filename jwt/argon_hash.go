package jwt

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/crypto/argon2"
)

var (
	ErrorInvalidArgonHashString   = errors.New("invalid argon hash string")
	ErrorIncompatibleArgonVersion = errors.New("incompatible argon version")
)

func PasswordHashArgon2i(password string, p model.PasswordHashArgonParams, salt, pepper []byte) string {
	// let's encode the hash into the standard text representation
	// https://github.com/P-H-C/phc-winner-argon2#command-line-utility
	// example: $argon2i$v=19$m=65536,t=2,p=4$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	slatAndPepper := append(salt, pepper...)
	hash := argon2.IDKey([]byte(password), slatAndPepper, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	salt64 := base64.RawStdEncoding.EncodeToString(salt)
	hash64 := base64.RawStdEncoding.EncodeToString(hash)

	enc := fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", string(model.PasswordHashArgon2i), argon2.Version, p.Memory, p.Iterations, p.Parallelism, salt64, hash64)
	return enc
}

func PasswordMatchArgon2i(password, hash string, pepper []byte) (bool, error) {
	p, salt, eh, err := decodeHashStringArgon2i(hash)
	if err != nil {
		return false, err
	}

	slatAndPepper := append(salt, pepper...)
	otherHash := argon2.IDKey([]byte(password), slatAndPepper, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// use subtle to prevent timing attacks
	if subtle.ConstantTimeCompare(eh, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHashStringArgon2i decode hash string and returns params, salt, hash and error.
func decodeHashStringArgon2i(hashStr string) (model.PasswordHashArgonParams, []byte, []byte, error) {
	params := model.PasswordHashArgonParams{}
	vals := strings.Split(hashStr, "$")
	if len(vals) != 6 {
		return params, nil, nil, ErrorInvalidArgonHashString
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return params, nil, nil, err
	}
	if version != argon2.Version {
		return params, nil, nil, ErrorIncompatibleArgonVersion
	}

	p := &params
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return params, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return params, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return params, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
