package api_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InviteRequest struct {
	Email       string `json:"email"`
	Role        string `json:"access_role"`
	CallbackURL string `json:"callback_url"`
}

func TestInvite(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)
	signature, _ := Signature(data, cfg.AppSecret)

	at := ""
	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			at = data["access_token"].(string)
			return nil
		})).
		JSONSchema("../../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()

	require.NotEmpty(t, at)
	data = fmt.Sprintf(`
	{ 
		"email": "%s",
		"access_role": "%s",
		"callback_url": "%s"
	}`, "invitee@madappgang.com", "admin", "http://localhost:3322")
	signature, _ = Signature(data, cfg.AppSecret)

	link := ""
	request.Post("/auth/invite").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Authorization", "Bearer "+at).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			link = data["link"].(string)
			return nil
		})).
		Type("json").
		Status(200).
		Done()

	assert.NotEmpty(t, link)
	assert.Contains(t, link, "email=invitee@madappgang.com")
	assert.Contains(t, link, `callbackUrl=http:%2F%2Flocalhost:3322`)
}
