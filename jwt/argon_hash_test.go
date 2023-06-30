package jwt_test

import (
	"encoding/base64"
	"testing"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestArgonHash(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashArgonParams
	salt := []byte("I am a salt!!!")
	pepper := []byte("I am a pepper!!!")

	hash := jwt.PasswordHashArgon2i(p, params, salt, pepper)

	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$argon2i$")
	assert.Contains(t, hash, "$v=")
	assert.Contains(t, hash, "$m=")
	assert.Contains(t, hash, ",t=")
	assert.Contains(t, hash, ",p=")
	assert.Contains(t, hash, base64.RawStdEncoding.EncodeToString(salt))
}

func TestArgonValidatePassword(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashArgonParams
	salt := []byte("I am a salt!!!")
	pepper := []byte("I am a pepper!!!")

	hash := jwt.PasswordHashArgon2i(p, params, salt, pepper)

	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$argon2i$")
	assert.Contains(t, hash, "$v=")

	match, err := jwt.PasswordMatchArgon2i("password111213", hash, pepper)
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = jwt.PasswordMatchArgon2i("!password111213", hash, pepper)
	assert.NoError(t, err)
	assert.False(t, match)

	match, err = jwt.PasswordMatchArgon2i("", hash, pepper)
	assert.NoError(t, err)
	assert.False(t, match)
}
