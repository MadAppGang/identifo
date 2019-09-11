package admin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

type file struct {
	Contents string `json:"contents"`
}

// GetStringifiedFile fetches static file from the static files storage,
// and returns its string representation.
func (ar *Router) GetStringifiedFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename, err := ar.parseStaticFileFullName(w, r)
		if err != nil {
			return
		}

		fileBytes, err := ar.staticFilesStorage.GetFile(filename)
		if err != nil {
			if err == model.ErrorNotFound {
				ar.Error(w, fmt.Errorf("File %s not found", filename), http.StatusNotFound, err.Error())
				return
			}
			ar.Error(w, err, http.StatusInternalServerError, err.Error())
			return
		}

		response := &file{Contents: string(fileBytes)}
		ar.ServeJSON(w, http.StatusOK, response)
	}
}

// UploadStringifiedFile uploads stringified static file to the storage.
func (ar *Router) UploadStringifiedFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename, err := ar.parseStaticFileFullName(w, r)
		if err != nil {
			return
		}

		f := new(file)
		if ar.mustParseJSON(w, r, f) != nil {
			return
		}

		if err := ar.staticFilesStorage.UploadFile(filename, []byte(f.Contents)); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, err.Error())
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

const oneMegabyte = int64(1 * 1024 * 1024)

// UploadADDAFile is for uploading Apple Developer Domain Association File.
func (ar *Router) UploadADDAFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, oneMegabyte)

		if err := r.ParseMultipartForm(oneMegabyte); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		formFile, _, err := r.FormFile("file")
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Cannot read file: %s", err.Error()))
			return
		}
		defer formFile.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, formFile); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot read file as bytes: %s", err.Error()))
			return
		}

		if err = ar.staticFilesStorage.UploadFile(model.AppleFilenames.DeveloperDomainAssociation, buf.Bytes()); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot upload file: %s", err.Error()))
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, oneMegabyte)

		if err := r.ParseMultipartForm(oneMegabyte); err != nil {
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

		if keys == nil {
			ar.Error(w, fmt.Errorf("Keys are empty"), http.StatusBadRequest, "")
			return
		}
		if keys.Private == nil {
			ar.Error(w, fmt.Errorf("Empty private key"), http.StatusBadRequest, "")
			return
		}
		if keys.Public == nil {
			ar.Error(w, fmt.Errorf("Empty public key"), http.StatusBadRequest, "")
			return
		}

		if err := ar.configurationStorage.InsertKeys(keys); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

func (ar *Router) parseStaticFileFullName(w http.ResponseWriter, r *http.Request) (string, error) {
	filename := r.URL.Query().Get("name")
	if len(filename) == 0 {
		err := fmt.Errorf("Empty filename")
		ar.Error(w, err, http.StatusBadRequest, "")
		return "", err
	}

	extension := r.URL.Query().Get("ext")
	if len(extension) != 0 {
		return strings.Join([]string{filename, extension}, "."), nil
	}
	return filename, nil
}
