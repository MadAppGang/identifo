package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHashBcrypt use pepper on top of default bcrypt encoding
// instead of bcrypt(password, cost)
// we do:
// hashedPassword = hmac(sha256(password), pepper)
// bcrypt(hashedPassword, cost)
func PasswordHashBcrypt(password string, params model.PasswordHashBcryptParams, salt, pepper []byte) (string, error) {
	passwordHmac := hmac.New(sha256.New, pepper) // we use sha256 to pepper the password
	passwordHmac.Write([]byte(password))
	bs, err := bcrypt.GenerateFromPassword(passwordHmac.Sum(nil), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(bs), nil
}

// PasswordMatchBcrypt check if password and hash matches
func PasswordMatchBcrypt(password, hash string, pepper []byte) (bool, error) {
	passwordHmac := hmac.New(sha256.New, pepper) // we use sha256 to pepper the password
	passwordHmac.Write([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hash), passwordHmac.Sum(nil))
	return err == nil, err
}
