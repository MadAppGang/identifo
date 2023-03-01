package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsAlwaysTrue(t *testing.T) {
	var message string
	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		Expect(t).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			err := data["error"].(map[string]interface{})
			message = err["message"].(string)
			return nil
		})).
		AssertFunc(dumpResponse).
		Status(500).
		Done()
	assert.NotEmpty(t, message)
	assert.Contains(t, message, "DefaultStorage settings: unsupported database type wrong_type")
}
