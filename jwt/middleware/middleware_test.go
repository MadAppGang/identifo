package middleware_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/madappgang/identifo/v2/jwt/middleware"
	"github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/model"
)

const (
	keyPath            = "../test_artifacts/public.pem"
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
		TokenType:      []string{middleware.TokenTypeAccess, middleware.TokenTypeRefresh},
		Audience:       []string{tokenAud, "ddddd"},
		Issuer:         []string{testIssuer, "dsfsd"},
	}

	wrongIssuerConfig := validator.Config{
		PubKeyEnvName:  "PK",
		PubKeyFileName: keyPath,
		TokenType:      []string{middleware.TokenTypeAccess},
		Audience:       []string{tokenAud},
		Issuer:         []string{"I am wrong issuer"},
	}

	specificUserIssuerConfig := validator.Config{
		PubKeyEnvName:  "PK",
		PubKeyFileName: keyPath,
		TokenType:      []string{middleware.TokenTypeAccess},
		Audience:       []string{tokenAud},
		Issuer:         []string{testIssuer},
		UserID:         []string{"user1"},
	}

	type args struct {
		c validator.Config
	}
	tests := []struct {
		name string
		args args
		want middleware.Error
	}{
		{"successfull get token", args{successConfig}, ""},
		{"invalid issuer", args{wrongIssuerConfig}, middleware.ErrorTokenIsInvalid},
		{"invalid user", args{specificUserIssuerConfig}, middleware.ErrorTokenIsInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meh := &mockErrorHandler{t: t, e: ""}

			handler, _ := middleware.JWT(meh, tt.args.c)

			req, _ := http.NewRequest(http.MethodGet, "testServer.URL", nil)
			req.Header.Add("Authorization", "Bearer "+tokenStringExample)

			next := func(rw http.ResponseWriter, r *http.Request) {
			}

			handler(nil, req, next)

			if tt.want != meh.e {
				t.Errorf("TestMiddleware() error = %T, wantErr %T %v", meh.e, tt.want, tt.want != meh.e)
			}
		})
	}
}

type mockErrorHandler struct {
	t *testing.T
	e middleware.Error
}

func (m *mockErrorHandler) Error(rw http.ResponseWriter, errorType middleware.Error, status int, description string) {
	m.e = errorType
}

func TestTokenFromContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want model.Token
	}{
		{
			name: "no token",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "nil token",
			ctx:  context.WithValue(context.Background(), model.TokenContextKey, nil),
			want: nil,
		},
		{
			name: "token exists",
			ctx:  context.WithValue(context.Background(), model.TokenContextKey, model.Token(&model.JWToken{})),
			want: &model.JWToken{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := middleware.TokenFromContext(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
