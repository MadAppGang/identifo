package runner_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/madappgang/identifo/test/runner"
)

func TestLoginAndRefreshToken(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, user1, user1Pswd)
	signature, _ := runner.Signature(data, appSecret)
	rt := ""

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", appID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_with_refresh_scheme.json").
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			rt = data["refresh_token"].(string)
			return nil
		})).
		Done()

	data = `{}`
	d := fmt.Sprintf("%d", time.Now().Unix())
	signature, _ = runner.Signature("/auth/token"+d, appSecret)

	request.Post("/auth/token").
		SetHeader("X-Identifo-ClientID", appID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+rt).
		SetHeader("X-Identifo-Timestamp", d).
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("data/jwt_token_with_refresh_scheme.json").
		Done()
}
