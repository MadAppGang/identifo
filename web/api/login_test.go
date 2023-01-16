package api_test

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/madappgang/identifo/v2/model"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
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
		JSONSchema("../../test/artifacts/api/jwt_token_with_refresh_scheme.json").
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
		JSONSchema("../../test/artifacts/api/jwt_token_scheme.json").
		Done()
}

// test wrong app ID login
func TestLoginWithWrongAppID(t *testing.T) {
	g := NewGomegaWithT(t)

	request.Post("/auth/login").
		SetHeader("X-Identifo-ClientID", "wrong_app_ID").
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data["error"]).To(MatchAllKeys(Keys{
				"id":               Equal("error.api.request.app_id.invalid"),
				"message":          Not(BeZero()),
				"detailed_message": Not(BeZero()),
				"status":           BeNumerically("==", 400),
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
		// AssertFunc(dumpResponse).
		Status(400).
		AssertFunc(validateJSON(func(data map[string]interface{}) error {
			g.Expect(data["error"]).To(MatchAllKeys(Keys{
				"id":      Equal("error.api.request.signature.invalid"),
				"message": Not(BeZero()),
				"status":  BeNumerically("==", 400),
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
			tokenStr = data["refresh_token"].(string)
			return nil
		})).
		Status(200).
		JSONSchema("../../test/artifacts/api/jwt_token_with_refresh_scheme.json").
		Done()

	tt, _ := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, nil)
	token := model.JWToken{JWT: tt}
	token.
}
