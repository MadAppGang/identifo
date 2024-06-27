package api_test

import (
	"fmt"
	"testing"
)

// ============================================================
// Login with phone number
// ============================================================

// test happy day login with phone number for one user
func TestLoginWithPhoneNumber(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"phone_number": "%s"
	}`, "+380123456789")
	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/request_phone_code").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		Done()

	data2 := fmt.Sprintf(`
	{
		"phone_number": "%s",
		"code": "%s",
		"scopes": [
			"offline",
			"smartrun"
		]
	}`, "+380123456789", "314159")
	signature, _ = Signature(data2, cfg.AppSecret)
	request.Post("/auth/phone_login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data2).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()
}

// test happy day login with phone number for two users
func TestLoginWithPhoneNumberTwoUsers(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"phone_number": "%s"
	}`, "+380123456781")
	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/request_phone_code").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		Done()

	data2 := fmt.Sprintf(`
	{
		"phone_number": "%s",
		"code": "%s",
		"scopes": [
			"offline",
			"smartrun"
		]
	}`, "+380123456781", "314159")
	signature, _ = Signature(data2, cfg.AppSecret)
	request.Post("/auth/phone_login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data2).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()

	data = fmt.Sprintf(`
	{
		"phone_number": "%s"
	}`, "+380123456782")
	signature, _ = Signature(data, cfg.AppSecret)

	request.Post("/auth/request_phone_code").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		Done()

	data2 = fmt.Sprintf(`
	{
		"phone_number": "%s",
		"code": "%s",
		"scopes": [
			"offline",
			"smartrun"
		]
	}`, "+380123456782", "314159")
	signature, _ = Signature(data2, cfg.AppSecret)
	request.Post("/auth/phone_login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data2).
		Expect(t).
		// AssertFunc(dumpResponse(t)).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()
}
