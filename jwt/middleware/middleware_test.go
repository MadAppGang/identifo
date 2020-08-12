package middleware_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/madappgang/identifo/jwt/middleware"
	"github.com/madappgang/identifo/jwt/validator"
)

const (
	keyPath            = "../public.pem"
	testIssuer         = "identifo.madappgang.com"
	tokenStringExample = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1OTcxNDY4OTIsImF1ZCI6InRlc3RfYXVkIiwiaXNzIjoiaWRlbnRpZm8ubWFkYXBwZ2FuZy5jb20iLCJ0eXBlIjoiYWNjZXNzIn0.BqdHOYtBPG9f7lZwPsV3OLNjd2y_vsSZlGCFbJOv2njaJ1poLBmw9VxthKU-L7Sr0X-E_yYldIGxV6ePryJuCg"
	tokenAud           = "test_aud"
	publicKeyString    = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAED3DoOWZMbqYc0OO1Ih628hB2Odhv
4mjl1vt0iBu3gTKz1XAk+YHG8aoI+42TJle6hawmTzrSD6khaNRaDQAKbg==
-----END PUBLIC KEY-----`
)

func TestMiddleware(t *testing.T) {

	os.Setenv("PK", publicKeyString)

	successConfig := validator.Config{
		PubKeyEnvName:  "PK",
		PubKeyFileName: keyPath,
		TokenType:      middleware.TokenTypeAccess,
		Audience:       tokenAud,
		Issuer:         testIssuer,
	}

	wrongIssuerConfig := validator.Config{
		PubKeyEnvName:  "PK",
		PubKeyFileName: keyPath,
		TokenType:      middleware.TokenTypeAccess,
		Audience:       tokenAud,
		Issuer:         "I am wrong issuer",
	}

	specificUserIssuerConfig := validator.Config{
		PubKeyEnvName:  "PK",
		PubKeyFileName: keyPath,
		TokenType:      middleware.TokenTypeAccess,
		Audience:       tokenAud,
		Issuer:         testIssuer,
		UserID:         "user1",
	}

	type args struct {
		c validator.Config
	}
	tests := []struct {
		name    string
		args    args
		want    middleware.Error
		wantErr bool
	}{
		{"successfull get token", args{successConfig}, "", false},
		{"invalid issuer", args{wrongIssuerConfig}, middleware.ErrorTokenIsInvalid, true},
		{"invalid user", args{specificUserIssuerConfig}, middleware.ErrorTokenIsInvalid, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meh := mockErrorHandler{t: t, e: ""}
			handler := middleware.Middleware(meh, tt.args.c)
			req, _ := http.NewRequest(http.MethodGet, "testServer.URL", nil)
			req.Header.Add("Authorization", "Bearer "+tokenStringExample)
			next := func(rw http.ResponseWriter, r *http.Request) {
			}
			handler(nil, req, next)
			if (tt.want != meh.e) != tt.wantErr {
				t.Errorf("TestMiddleware() error = %v, wantErr %v", meh.e, tt.wantErr)
			}
		})
	}
}

type mockErrorHandler struct {
	t *testing.T
	e middleware.Error
}

func (m mockErrorHandler) Error(rw http.ResponseWriter, errorType middleware.Error, status int, description string) {
	m.e = errorType
}
