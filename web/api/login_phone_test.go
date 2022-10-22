package api_test

import (
	"fmt"
	"testing"
)

// ============================================================
// Login with phone number
// ============================================================

// test happy day login
func TestLoginWithPhoneNumber(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)
	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()
}
