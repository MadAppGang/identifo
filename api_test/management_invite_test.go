package api_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	sig "github.com/madappgang/digestsig"
	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagementResetToken(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"

	}`, cfg.User1)

	tokenString := ""
	u, _ := url.Parse(cfg.ServerURL)
	sd := sig.SigningData{
		Method:      "POST",
		BodyMD5:     sig.GetMD5([]byte(data)),
		ContentType: "application/json",
		Date:        time.Now().Format(time.RFC3339),
		Expires:     time.Now().Add(time.Hour).Unix(),
		Host:        u.Host,
	}
	fmt.Println(sd.String())
	signature := sig.SignString(sd.String(), []byte(cfg.ManagementKeySecret1))

	request.Post("/management/token/reset_password").
		SetHeader("Content-Type", sd.ContentType).
		SetHeader("Expires", fmt.Sprintf("%d", sd.Expires)).
		SetHeader("Date", sd.Date).
		SetHeader("Content-MD5", sd.BodyMD5).
		SetHeader("Digest", fmt.Sprintf("%s%s", sig.DigestHeaderSHAPrefix, signature)).
		SetHeader(sig.KeyIDHeaderKey, cfg.ManagementKeyID1).
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			tokenString = data["token"].(string)
			return nil
		})).
		Done()

	require.NotEmpty(t, tokenString)

	token, err := jwt.ParseTokenString(tokenString)
	require.NoError(t, err)
	require.NotNil(t, token)

	require.Equal(t, string(model.TokenTypeReset), token.Type())
	assert.Equal(t, "identifo", token.Audience())
}

func TestManagementInactiveKey(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"

	}`, cfg.User1)

	u, _ := url.Parse(cfg.ServerURL)
	sd := sig.SigningData{
		Method:      "POST",
		BodyMD5:     sig.GetMD5([]byte(data)),
		ContentType: "application/json",
		Date:        time.Now().Format(time.RFC3339),
		Expires:     time.Now().Add(time.Hour).Unix(),
		Host:        u.Host,
	}
	fmt.Println(sd.String())
	signature := sig.SignString(sd.String(), []byte(cfg.ManagementKeySecret2))

	body := ""

	request.Post("/management/token/reset_password").
		SetHeader("Content-Type", sd.ContentType).
		SetHeader("Expires", fmt.Sprintf("%d", sd.Expires)).
		SetHeader("Date", sd.Date).
		SetHeader("Content-MD5", sd.BodyMD5).
		SetHeader("Digest", fmt.Sprintf("%s%s", sig.DigestHeaderSHAPrefix, signature)).
		SetHeader(sig.KeyIDHeaderKey, cfg.ManagementKeyID2).
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Status(400).
		AssertFunc(validateBodyText(func(b string) error {
			body = b
			return nil
		})).
		Done()

	require.NotEmpty(t, body)
	assert.Contains(t, body, "error.native.login.ma.key.inactive")
}

func TestManagementInviteNoKeyID(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"

	}`, cfg.User1)

	body := ""
	request.Post("/management/reset_password_token").
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Status(400).
		AssertFunc(validateBodyText(func(b string) error {
			body = b
			return nil
		})).
		Done()

	require.NotEmpty(t, body)
	assert.Contains(t, body, "error.native.login.ma.no.key.id")
}

func TestManagementInviteWrongKeyID(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"email": "%s"

	}`, cfg.User1)

	body := ""
	request.Post("/management/reset_password_token").
		SetHeader("Content-Type", "application/json").
		SetHeader("Digest", "AABBCCDDSS").
		SetHeader("X-Nl-Key-Id", "AABBCCDDSS").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Status(400).
		AssertFunc(validateBodyText(func(b string) error {
			body = b
			return nil
		})).
		Done()

	require.NotEmpty(t, body)
	assert.Contains(t, body, "error.native.login.ma.error.key.with.id")
}
