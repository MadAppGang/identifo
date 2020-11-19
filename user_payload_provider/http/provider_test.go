package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	hu "github.com/madappgang/identifo/user_payload_provider/http"
)

const (
	secret  = "super_secret"
	appId   = "12345"
	appName = "I am the web app"
	userId  = "09876543d21"
)

func Test_provider_UserPayloadForApp(t *testing.T) {
	//precalculated value from https://www.freeformatter.com/hmac-generator.html#ad-output
	//using input: {"app_id":"12345","app_name":"I am the web app","user_id":"09876543d21"}
	expectedDigest := "b9a6be00d9656fee55165749596a321a2c33abf795d61d5714c44715b81371a0"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		digest := r.Header["Digest"][0][len("SHA-256="):]
		fmt.Println(digest)
		if digest != expectedDigest {
			t.Errorf("wrong digest %v, expected %v", digest, expectedDigest)
		}
		fmt.Fprintln(w, `{"city" : "Sydney"}`)
	}))
	defer ts.Close()

	p, err := hu.NewUserPayloadProvider(secret, ts.URL)
	if err != nil {
		t.Errorf("unable to create provider with error %v", err)
		t.FailNow()
	}

	if p == nil {
		t.Error("provider should not be empty")
		t.FailNow()
	}

	payload, err := p.UserPayloadForApp(appId, appName, userId)
	if err != nil {
		t.Errorf("unable to get data payload with error %v", err)
		t.FailNow()
	}

	if payload["city"] != "Sydney" {
		t.Errorf("Got unexpected value  %v for key \"city\", expected \"Sydney\"", payload["city"])
		t.FailNow()
	}

}
