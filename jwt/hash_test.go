package jwt_test

import (
	"encoding/base64"
	"testing"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashParams
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHash(p, params, pepper)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$argon2i$")
	assert.Contains(t, hash, "$v=")
	assert.Contains(t, hash, "$m=")
	assert.Contains(t, hash, ",t=")
	assert.Contains(t, hash, ",p=")
}

func TestOtherHash(t *testing.T) {
	p := "password111213"
	params := model.PasswordHashParams{
		Type:       model.PasswordHashBcrypt,
		SaltLength: 16,
		Bcrypt:     &model.DefaultPasswordHashBcryptParams,
	}
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHash(p, params, pepper)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$2a$")
	assert.NotContains(t, hash, base64.RawStdEncoding.EncodeToString(pepper))
}

func TestInvalidHashAlg(t *testing.T) {
	p := "password111213"
	params := model.PasswordHashParams{
		Type:       "whatever",
		SaltLength: 16,
		Bcrypt:     &model.DefaultPasswordHashBcryptParams,
	}
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHash(p, params, pepper)

	assert.Error(t, err)
	assert.Equal(t, err, jwt.ErrorUnknownPasswordHashType)
	assert.Empty(t, hash)
}

func TestInvalidHashParams(t *testing.T) {
	p := "password111213"
	params := model.PasswordHashParams{
		Type:       model.PasswordHashArgon2i,
		SaltLength: 16,
	}
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHash(p, params, pepper)

	assert.Error(t, err)
	assert.Equal(t, err, jwt.ErrorHashParamsMissing)
	assert.Empty(t, hash)
}

func TestPasswordMatch(t *testing.T) {
	p := "password111213"
	params := model.DefaultPasswordHashParams
	pepper := []byte("I am a pepper!!!")

	hash, err := jwt.PasswordHash(p, params, pepper)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	match, err := jwt.PasswordMatch(p, hash, pepper)
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = jwt.PasswordMatch("!"+p, hash, pepper)
	assert.NoError(t, err)
	assert.False(t, match)

	match, err = jwt.PasswordMatch(p, hash, append(pepper, []byte("!")...))
	assert.NoError(t, err)
	assert.False(t, match)
}
