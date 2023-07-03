package jwt

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/madappgang/identifo/v2/model"
)

var (
	ErrorHashParamsMissing       = errors.New("missing password hash algorithm parameters")
	ErrorUnknownPasswordHashType = errors.New("unsupported password hash algorithm")
)

// PasswordHash generates a hash of a password using the given parameters and secret sal (pepper).
func PasswordHash(password string, params model.PasswordHashParams, pepper []byte) (string, error) {
	salt, err := getRandomSalt(params.SaltLength)
	if err != nil {
		return "", err
	}

	if params.Type == model.PasswordHashDefault {
		params = model.DefaultPasswordHashParams
	}

	var hash string
	switch params.Type {
	case model.PasswordHashArgon2i:
		if params.Argon == nil {
			err = ErrorHashParamsMissing
		} else {
			hash = PasswordHashArgon2i(password, *params.Argon, salt, pepper)
		}
	case model.PasswordHashBcrypt:
		if params.Bcrypt == nil {
			err = ErrorHashParamsMissing
		} else {
			hash, err = PasswordHashBcrypt(password, *params.Bcrypt, pepper)
		}
	default:
		err = ErrorUnknownPasswordHashType
	}

	return hash, err
}

const bcryptHashPrefix = "$2" // it could be $2a, $2b, $2y, $2x

func PasswordMatch(password, hash string, pepper []byte) (bool, error) {
	// first, we need to extract the algorithm type
	if strings.HasPrefix(hash, bcryptHashPrefix) {
		return PasswordMatchBcrypt(password, hash, pepper)
	} else if strings.HasPrefix(hash, fmt.Sprintf("$%s$", string(model.PasswordHashArgon2i))) {
		return PasswordMatchArgon2i(password, hash, pepper)
	} else {
		return false, ErrorUnknownPasswordHashType
	}
}

func getRandomSalt(n uint32) ([]byte, error) {
	unencodedSalt := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, unencodedSalt)
	if err != nil {
		return nil, err
	}

	return unencodedSalt, err
}
