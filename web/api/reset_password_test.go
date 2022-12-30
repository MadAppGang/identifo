package api_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetPasswordWithCustomURL(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s",
		"reset_page_url": "%s"
	}`, cfg.User1, "https://customurl.com")
	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/request_reset_password").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../../test/artifacts/api/response_ok.json").
		Done()

	// if running local server (the email sever will not be nil then), check the email content
	if emailService != nil {
		messages := emailService.Messages()
		require.GreaterOrEqual(t, len(messages), 1) // at least one message should be send
		lastMessage := messages[len(messages)-1]
		assert.Contains(t, lastMessage, cfg.User1)
		assert.Contains(t, lastMessage, "href=\"https://customurl.com?")
		assert.Contains(t, lastMessage, fmt.Sprintf("appId=%s", cfg.AppID))
		fmt.Printf("\nEmail:\n%s\n", lastMessage)
	}
}

func TestResetPasswordWithAppSpecificURL(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"
	}`, cfg.User1)
	signature, _ := Signature(data, cfg.AppSecret2)

	request.Post("/auth/request_reset_password").
		SetHeader("X-Identifo-ClientID", cfg.AppID2).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../../test/artifacts/api/response_ok.json").
		Done()

	// if running local server (the email sever will not be nil then), check the email content
	if emailService != nil {
		messages := emailService.Messages()
		require.GreaterOrEqual(t, len(messages), 1) // at least one message should be send
		lastMessage := messages[len(messages)-1]
		assert.Contains(t, lastMessage, cfg.User1)
		assert.Contains(t, lastMessage, "href=\"http://rewrite.com/login/cusom?")
		assert.Contains(t, lastMessage, fmt.Sprintf("appId=%s", cfg.AppID2))
	}
}

func TestResetPasswordWithDefaultURL(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"
	}`, cfg.User1)
	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/request_reset_password").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../../test/artifacts/api/response_ok.json").
		Done()

	// if running local server (the email sever will not be nil then), check the email content
	if emailService != nil {
		messages := emailService.Messages()
		require.GreaterOrEqual(t, len(messages), 1) // at least one message should be send
		lastMessage := messages[len(messages)-1]
		assert.Contains(t, lastMessage, cfg.User1)
		assert.Contains(t, lastMessage, "href=\"http://localhost:8081/web/password/reset?")
		assert.Contains(t, lastMessage, fmt.Sprintf("appId=%s", cfg.AppID))
	}
}
