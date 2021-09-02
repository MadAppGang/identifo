package admin

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
)

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, oneMegabyte)

		if err := r.ParseMultipartForm(oneMegabyte); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		keys := model.JWTKeys{}

		/// private key read
		_, prHeader, err := r.FormFile("private")
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("error parsing parsing private file: %s", err.Error()))
			return
		}
		fp, err := prHeader.Open()
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error getting private key: %s", err.Error()))
			return
		}
		defer fp.Close()
		keys.Private = fp

		/// public key read
		_, pubHeader, err := r.FormFile("public")
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("error parsing parsing public file: %s", err.Error()))
			return
		}
		fpub, err := pubHeader.Open()
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error getting public key: %s", err.Error()))
			return
		}
		defer fpub.Close()
		keys.Public = fpub

		if err := ar.server.Storages().Key.ReplaceKeys(keys); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// GetJWTKeys returns public and private JWT keys currently used by Identifo
func (ar *Router) GetJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := ar.server.Storages().Key.GetKeys()
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, keys)
	}
}
