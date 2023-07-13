package api_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// Login with username and password
// ============================================================

// test happy day login
func TestLogin(t *testing.T) {
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
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()
}

// test happy day login, with no refresh token included
func TestLoginWithNoRefresh(t *testing.T) {
	g := NewGomegaWithT(t)

	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": []
	}`, cfg.User1, cfg.User1Pswd)

	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data).To(MatchKeys(IgnoreExtras|IgnoreMissing, Keys{
				"access_token":  Not(BeZero()),
				"refresh_token": BeZero(),
			}))
			return nil
		})).
		Type("json").
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_scheme.json").
		Done()
}

// test wrong app ID login
func TestLoginWithWrongAppID(t *testing.T) {
	g := NewGomegaWithT(t)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "wrong_app_ID").
		Expect(t).
		AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data["error"]).To(MatchAllKeys(Keys{
				"id":       Equal(string(l.ErrorStorageAPPFindByIDError)),
				"message":  Not(BeZero()),
				"location": Not(BeZero()),
				"status":   BeNumerically("==", 400),
			}))
			return nil
		})).
		Type("json").
		Status(400).
		Done()
}

// test wrong signature for mobile app
func TestLoginWithWrongSignature(t *testing.T) {
	g := NewGomegaWithT(t)

	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)

	signature, _ := Signature(data, cfg.AppSecret)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature+"_wrong").
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		AssertFunc(dumpResponse).
		Status(400).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data["error"]).To(MatchAllKeys(Keys{
				"id":       Equal(string(l.ErrorAPIRequestSignatureInvalid)),
				"message":  Not(BeZero()),
				"status":   BeNumerically("==", 400),
				"location": Not(BeZero()),
			}))
			return nil
		})).
		Done()
}

func TestLoginTokenClaims(t *testing.T) {
	data := fmt.Sprintf(`
	{
		"username": "%s",
		"password": "%s",
		"scopes": ["offline"]
	}`, cfg.User1, cfg.User1Pswd)
	signature, _ := Signature(data, cfg.AppSecret)
	tokenStr := ""

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		BodyString(data).
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			tokenStr = data["access_token"].(string)
			return nil
		})).
		Status(200).
		JSONSchema("../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()

	body := ""
	// Get public key
	request.Get("/.well-known/jwks.json").
		Expect(t).
		Type("json").
		AssertFunc(validateBodyText(func(b string) error {
			body = b
			return nil
		})).
		Status(200).
		Done()

	jwks, err := keyfunc.NewJSON(json.RawMessage(body))
	assert.NoError(t, err)

	tt, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, jwks.Keyfunc)
	assert.NoError(t, err)
	assert.Equal(t, "ES256", tt.Method.Alg())
	assert.True(t, tt.Valid)

	token := model.JWToken{Token: *tt}
	assert.False(t, token.New)
	assert.Equal(t, 1, len(token.Payload()))
	assert.Equal(t, "59fd884d8f6b180001f5b4e2", token.Audience())
	assert.Equal(t, "63c6273a46504e3abdc00fc6", token.Subject())
	assert.Equal(t, "http://localhost", token.Issuer())
	assert.Equal(t, "access", token.Type())

	assert.GreaterOrEqual(t, len(token.Audience()), 1)
}
