package runner_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/test/runner"
)

// test register with email and password
func TestRegisterWithEmail(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, user2, user2Pswd)

	signature, _ := runner.Signature(data, appSecret)

	request.Post("/auth/register").
		SetHeader("X-Identifo-ClientID", appID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_with_refresh_scheme.json").
		Done()
}

// test register and logout with access token
func TestRegisterWithEmailAndLogout(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, user3, user3Pswd)

	signature, _ := runner.Signature(data, appSecret)

	at := ""
	rt := ""

	request.Post("/auth/register").
		SetHeader("X-Identifo-ClientID", appID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			at = data["access_token"].(string)
			rt = data["refresh_token"].(string)
			return nil
		})).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_with_refresh_scheme.json").
		Done()

	logoutData := fmt.Sprintf(`
	{
		"refresh_token": "%s"
	}`, rt)

	signatureLogout, _ := runner.Signature(logoutData, appSecret)
	request.Post("/me/logout").
		SetHeader("X-Identifo-ClientID", appID).
		SetHeader("Digest", "SHA-256="+signatureLogout).
		SetHeader("Authorization", "Bearer "+at).
		SetHeader("Content-Type", "application/json").
		BodyString(logoutData).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		Done()
}
