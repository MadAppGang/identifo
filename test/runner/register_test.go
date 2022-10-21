package runner_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/test/runner"
)

// test register with email and password
func TestRegisterWithEmail(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline", "smartrun"]
	}`, cfg.User2, cfg.User2Pswd)

	signature, _ := runner.Signature(data, cfg.AppSecret)

	request.Post("/auth/register").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
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
		"scopes": ["offline", "smartrun"]
	}`, cfg.User3, cfg.User3Pswd)

	signature, _ := runner.Signature(data, cfg.AppSecret)

	at := ""
	rt := ""

	request.Post("/auth/register").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
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

	signatureLogout, _ := runner.Signature(logoutData, cfg.AppSecret)
	request.Post("/me/logout").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
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
