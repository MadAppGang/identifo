package http

import (
	"net/http"

	jwk "github.com/mendsley/gojwk"
)

func (ar *apiRouter) Configuration() http.HandlerFunc {

	type configurationResponse struct {
		Issuer                           string   `json:"issuer"`
		AuthorizationEndpoint            string   `json:"authorization_endpoint"`
		TokenEndpoint                    string   `json:"token_endpoint,omitempty"`
		JwksURI                          string   `json:"jwks_uri"`
		ResponseTypesSupported           []string `json:"response_types_supported"`
		SubjectTypesSupported            []string `json:"subject_types_supported"`
		IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	}

	iss := ar.tokenService.Issuer()
	configuration := configurationResponse{
		Issuer:                           iss,
		AuthorizationEndpoint:            iss + "/auth/login",
		TokenEndpoint:                    iss + "/auth/token",
		JwksURI:                          iss + "/.well-known/jwks",
		ResponseTypesSupported:           []string{"id_token", "token"},
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"ES256"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, configuration)
	}
}

func (ar *apiRouter) ServeJWKS() http.HandlerFunc {
	type JWKS struct {
		Keys []*jwk.Key `json:"keys"`
	}

	jwks := JWKS{
		Keys: []*jwk.Key{ar.tokenService.JWK()},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, jwks)
	}
}
