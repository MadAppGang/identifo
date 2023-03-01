package api_test

import (
	"fmt"
	"testing"
	"time"
)

func TestLoginAndRefreshToken(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)
	signature, _ := Signature(data, cfg.AppSecret)
	rt := ""

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			rt = data["refresh_token"].(string)
			return nil
		})).
		Done()

	d := fmt.Sprintf("%d", time.Now().Unix())
	signature, _ = Signature("/auth/token"+d, cfg.AppSecret)

	request.Post("/auth/token").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+rt).
		SetHeader("X-Identifo-Timestamp", d).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_refresh_token.json").
		Done()
}

func TestLoginAndRefreshTokenWithNewRefresh(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)
	signature, _ := Signature(data, cfg.AppSecret)
	rt := ""

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			rt = data["refresh_token"].(string)
			return nil
		})).
		Done()

	data = `{ "scopes": ["offline"] }`
	signature, _ = Signature(data, cfg.AppSecret)

	request.Post("/auth/token").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+rt).
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_refresh_token_with_new_refresh.json").
		Done()
}
