package admin

import (
	"net/http"
	"strings"

	jwt "github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

type keys struct {
	Private   string `json:"private,omitempty"`
	Public    string `json:"public,omitempty"`
	Algorithm string `json:"alg,omitempty"`
}

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		k := keys{}

		if err := ar.mustParseJSON(w, r, &k); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIJsonParseError, err.Error())
			return
		}

		if _, _, err := jwt.LoadPrivateKeyFromPEMString(k.Private); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelPrivateKeyEncoding, err.Error())
			return
		}

		if err := ar.server.Storages().Key.ReplaceKey([]byte(k.Private)); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeySave, err.Error())
			return
		}

		key, err := ar.server.Storages().Key.LoadPrivateKey()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeyLoad, err.Error())
			return
		}

		ar.server.Services().Token.SetPrivateKey(key)

		newkeys := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeyEncode, err.Error())
			return
		}
		newkeys.Public = publicPEM
		newkeys.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, locale, http.StatusOK, newkeys)
	}
}

// GetJWTKeys returns public and private JWT keys currently used by Identifo
func (ar *Router) GetJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		k := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPublicKeyEncode, err)
			return
		}
		k.Public = publicPEM

		private, ok := r.URL.Query()["include_private_key"]
		if ok && len(private) > 0 && strings.ToUpper(private[0]) == "TRUE" {
			private := ar.server.Services().Token.PrivateKey()
			privatePEM, err := jwt.MarshalPrivateKeyToPEM(private)
			if err != nil {
				ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeyEncode, err)
				return
			}
			k.Private = privatePEM
		}
		k.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, locale, http.StatusOK, k)
	}
}

// GenerateNewSecret generate new secret key, save it and return new public key
func (ar *Router) GenerateNewSecret() http.HandlerFunc {
	type payload struct {
		Alg string `json:"alg,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		p := payload{}

		if err := ar.mustParseJSON(w, r, &p); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIJsonParseError, err.Error())
			return
		}

		var alg model.TokenSignatureAlgorithm
		switch strings.ToLower(p.Alg) {
		case "es256":
			alg = model.TokenSignatureAlgorithmES256
		case "rs256":
			alg = model.TokenSignatureAlgorithmRS256
		default:
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelKeyAlgUnsupported, p.Alg)
			return
		}

		privateKey, err := jwt.GenerateNewPrivateKey(alg)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
			return
		}

		privateKeyPEM, err := jwt.MarshalPrivateKeyToPEM(privateKey)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeyEncode, err)
			return
		}

		if err := ar.server.Storages().Key.ReplaceKey([]byte(privateKeyPEM)); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeySave, err)
			return
		}

		key, err := ar.server.Storages().Key.LoadPrivateKey()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPrivateKeyLoad, err)
			return
		}

		ar.server.Services().Token.SetPrivateKey(key)

		newkeys := keys{}
		public := ar.server.Services().Token.PublicKey()
		publicPEM, err := jwt.MarshalPublicKeyToPEM(public)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelPublicKeyEncode, err)
			return
		}
		newkeys.Public = publicPEM
		newkeys.Algorithm = ar.server.Services().Token.Algorithm()

		ar.ServeJSON(w, locale, http.StatusOK, newkeys)
	}
}
