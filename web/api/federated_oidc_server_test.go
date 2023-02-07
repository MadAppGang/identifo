package api_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	ijwt "github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/api"
)

type OIDCCfg struct {
	api.OIDCConfiguration
	AuthURL  string `json:"authorization_endpoint"`
	TokenURL string `json:"token_endpoint"`
}

// testOIDCServer creates fake OIDC provider server for testing purposes.
func testOIDCServer() (*httptest.Server, context.CancelFunc) {
	privateKey, alg, err := ijwt.LoadPrivateKeyFromPEMFile("../../jwt/test_artifacts/private.pem")
	if err != nil {
		panic(err)
	}

	pk := privateKey.(*ecdsa.PrivateKey)
	publicKey := pk.Public()

	algName := strings.ToUpper(alg.String())

	cfg := &OIDCCfg{
		OIDCConfiguration: api.OIDCConfiguration{
			ScopesSupported:        []string{"idtoken"},
			SupportedIDSigningAlgs: []string{algName},
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(cfg)
	})

	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		jwk := api.CreateJWK(algName, "kid", publicKey)
		result := map[string]interface{}{"keys": []interface{}{jwk}}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := r.URL.Query().Get("redirect_uri")
		state := r.URL.Query().Get("state")
		code := "test_code"

		redirectUrl = redirectUrl + "?code=" + code + "&state=" + state

		http.Redirect(w, r, redirectUrl, http.StatusFound)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		idt, err := model.NewTokenWithClaims(jwt.SigningMethodES256, "kid", jwt.MapClaims{
			"sub":    "abc",
			"emails": []string{"some@example.com"},
			"email":  "some@example.com",
			"iss":    cfg.Issuer,
			"aud":    "test",
			"exp":    time.Now().Add(time.Hour).Unix(),
			"iat":    time.Now().Unix(),
		}).SignedString(privateKey)
		if err != nil {
			panic(err)
		}

		jt := struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"` // at least PayPal returns string, while most return number
			IDToken      string `json:"id_token"`
		}{
			IDToken:      idt,
			AccessToken:  "test_token",
			RefreshToken: "test_token",
			TokenType:    "Bearer",
			ExpiresIn:    12000,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jt)
	})

	oidcServer := httptest.NewServer(mux)

	cfg.Issuer = oidcServer.URL
	cfg.AuthURL = oidcServer.URL + "/auth"
	cfg.TokenURL = oidcServer.URL + "/token"
	cfg.JwksURI = oidcServer.URL + "/.well-known/jwks.json"

	return oidcServer, func() { oidcServer.Close() }
}
