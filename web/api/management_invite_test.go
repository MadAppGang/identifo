package api_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagementInvite(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"

	}`, cfg.User1)

	tokenString := ""
	request.Post("/management/reset_password_token").
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			tokenString = data["token"].(string)
			return nil
		})).
		// JSONSchema("../../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()

	require.NotEmpty(t, tokenString)

	token, err := jwt.ParseTokenString(tokenString)
	require.NoError(t, err)
	require.NotNil(t, token)

	require.Equal(t, string(model.TokenTypeReset), token.Type())
	assert.Equal(t, "identifo", token.Audience())
}
