package jwt_test

import (
	"encoding/base64"
	"testing"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestBcryptHash(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashBcryptParams
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHashBcrypt(p, params, pepper)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$2a$")
	assert.NotContains(t, hash, base64.RawStdEncoding.EncodeToString(pepper))
}

func TestBCryptValidatePassword(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashBcryptParams
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHashBcrypt(p, params, pepper)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	match, err := jwt.PasswordMatchBcrypt("password111213", hash, pepper)
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = jwt.PasswordMatchBcrypt("password111213", hash, append(pepper, []byte("1")...))
	assert.NoError(t, err)
	assert.False(t, match)

	match, err = jwt.PasswordMatchBcrypt("!password111213", hash, pepper)
	assert.NoError(t, err)
	assert.False(t, match)
}
