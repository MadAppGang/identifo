package admin

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/madappgang/identifo/jwt"
)

type keys struct {
	Private string `json:"private,omitempty"`
	Public  string `json:"public,omitempty"`
}

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := keys{}

		if err := ar.mustParseJSON(w, r, &k); err != nil {
			ar.Error(w, fmt.Errorf("error parsing keys: %v", err), http.StatusBadRequest, "")
			return
		}

		if _, _, err := jwt.LoadPrivateKeyFromString(k.Private); err != nil {
			ar.Error(w, fmt.Errorf("error decoding private key: %v", err), http.StatusBadRequest, "")
			return
		}

		if err := ar.server.Storages().Key.ReplaceKey([]byte(k.Private)); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		key, err := ar.server.Storages().Key.LoadPrivateKey()
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		ar.server.Services().Token.SetPrivateKey(key)
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// GetJWTKeys returns public and private JWT keys currently used by Identifo
func (ar *Router) GetJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		k.Public = publicPEM

		private, ok := r.URL.Query()["include_private_key"]
		if ok && len(private) > 0 && strings.ToUpper(private[0]) == "TRUE" {
			private := ar.server.Services().Token.PrivateKey()
			privatePEM, err := jwt.MarshalPrivateKeyToPEM(private)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
			k.Private = privatePEM
		}

		ar.ServeJSON(w, http.StatusOK, k)
	}
}
