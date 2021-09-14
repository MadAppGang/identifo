package admin

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

type keys struct {
	Private   string `json:"private,omitempty"`
	Public    string `json:"public,omitempty"`
	Algorithm string `json:"alg,omitempty"`
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

		newkeys := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		newkeys.Public = publicPEM
		newkeys.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, http.StatusOK, newkeys)
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
		k.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, http.StatusOK, k)
	}
}

// GenerateNewSecret generate new secret key, save it and return new public key
func (ar *Router) GenerateNewSecret() http.HandlerFunc {
	type payload struct {
		Alg string `json:"alg,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		p := payload{}

		if err := ar.mustParseJSON(w, r, &p); err != nil {
			ar.Error(w, fmt.Errorf("error parsing keys: %v", err), http.StatusBadRequest, "")
			return
		}

		var alg model.TokenSignatureAlgorithm
		switch strings.ToLower(p.Alg) {
		case "es256":
			alg = model.TokenSignatureAlgorithmES256
		case "rs256":
			alg = model.TokenSignatureAlgorithmRS256
		default:
			ar.Error(w, fmt.Errorf("unsupported algorithm in payload: %s", p.Alg), http.StatusBadRequest, "")
			return
		}

		privateKey, err := jwt.GenerateNewPrivateKey(alg)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		privateKeyPEM, err := jwt.MarshalPrivateKeyToPEM(privateKey)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		if err := ar.server.Storages().Key.ReplaceKey([]byte(privateKeyPEM)); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		key, err := ar.server.Storages().Key.LoadPrivateKey()
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		ar.server.Services().Token.SetPrivateKey(key)

		newkeys := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		newkeys.Public = publicPEM
		newkeys.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, http.StatusOK, newkeys)
	}
}
