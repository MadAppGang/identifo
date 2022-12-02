package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppSettings(t *testing.T) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	// GET requests does not have body,
	// that is why the signature should be based on URL and timestamp
	signature, _ := Signature("/auth/app_settings"+ts, cfg.AppSecret)

	request.Get("/auth/app_settings").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("X-Identifo-Timestamp", ts).
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		Expect(t).
		// AssertFunc(dumpResponse).
		Type("json").
		Status(200).
		JSONSchema("../../test/artifacts/api/app_settings_scheme.json").
		Done()
}

func TestAppSettingsCORS(t *testing.T) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	// GET requests does not have body,
	// that is why the signature should be based on URL and timestamp
	signature, _ := Signature("/auth/app_settings"+ts, cfg.AppSecret)
	var header http.Header

	request.Request().
		Method("OPTIONS").
		Path("/auth/app_settings").
		SetHeader("X-Identifo-ClientID", cfg.AppID).
		SetHeader("X-Identifo-Timestamp", ts).
		SetHeader("Origin", "http://localhost:3000").
		SetHeader("Access-Control-Request-Method", "GET").
		SetHeader("Digest", "SHA-256="+signature).
		SetHeader("Content-Type", "application/json").
		Expect(t).
		// AssertFunc(dumpResponse).
		AssertFunc(func(rsp *http.Response, r *http.Request) error {
			header = rsp.Header
			return nil
		}).
		Status(204).
		Done()

	fmt.Printf("Header %+v\n", header)
	assert.Contains(t, header.Get("Access-Control-Allow-Origin"), "http://localhost:3000")
	assert.Contains(t, header.Get("Access-Control-Allow-Methods"), "GET")
}


