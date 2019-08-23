package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/madappgang/identifo/model"
)

// UploadAppleDomainAssociation is for uploading Apple Domain Association File.
func (ar *Router) UploadAppleDomainAssociation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024 * 1); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		formFile, _, err := r.FormFile("file")
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Cannot read file: %s", err.Error()))
			return
		}
		defer formFile.Close()

		file, err := os.OpenFile(ar.ServerSettings.StaticFiles.AppleDomainAssociation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot open file: %s", err.Error()))
			return
		}
		defer file.Close()

		if _, err = io.Copy(file, formFile); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot save file: %s", err.Error()))
			return
		}

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024 * 1); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		formKeys := r.MultipartForm.File["keys"]

		keys := &model.JWTKeys{}

		for _, fileHeader := range formKeys {
			f, err := fileHeader.Open()
			if err != nil {
				ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error uploading key: %s", err.Error()))
				return
			}
			defer f.Close()

			switch fileHeader.Filename {
			case "private.pem":
				keys.Private = f
			case "public.pem":
				keys.Public = f
			default:
				ar.Error(w, fmt.Errorf("Invalid key field name '%s'", fileHeader.Filename), http.StatusBadRequest, "")
				return
			}
		}

		if err := ar.configurationStorage.InsertKeys(keys); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
