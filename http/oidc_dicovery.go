package http

import "net/http"

//https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
//Identifo is not OIDC provider, that's why we are providing here information only for token validation
type OIDCConfiguration struct {
	Issuer                 string   `json:"issuer"`
	JwksUri                string   `json:"jwks_uri"`
	ScopesSupported        []string `json:"scopes_supported"`
	SupportedIDSigningAlgs []string `json:"id_token_signing_alg_values_supported"`
}

//OIDCConfiguration - is endpoint to provide  OpenID Connect Discovery information (https://openid.net/specs/openid-connect-discovery-1_0.html)
//it should return  RFC5785 compatible documentation (https://tools.ietf.org/html/rfc5785)
//this will allow ise Identifo as Federated identity provider
//for example AWS AppSync (https://docs.aws.amazon.com/appsync/latest/devguide/security.html#openid-connect-authorization)
func (ar *apiRouter) OIDCConfiguration() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if ar.oidcConfiguration == nil {
			ar.oidcConfiguration = &OIDCConfiguration{
				Issuer:                 ar.tokenService.Issuer(),
				JwksUri:                ar.tokenService.Issuer() + "/.well-known/jwks.json",
				ScopesSupported:        ar.userStorage.Scopes(),
				SupportedIDSigningAlgs: []string{ar.tokenService.Algorithm()},
			}
		}
		ar.ServeJSON(w, http.StatusOK, ar.oidcConfiguration)
	}
}
