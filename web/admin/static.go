package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

// GetStaticFile fetches static file from the static files storage.
func (ar *Router) GetStaticFile() http.HandlerFunc {
	type file struct {
		Contents string `json:"contents"`
	}
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

func (ar *Router) parseStaticFileFullName(w http.ResponseWriter, r *http.Request) (filepath string, err error) {
	filename := r.URL.Query().Get("name")
	if len(filename) == 0 {
		err = fmt.Errorf("Empty filename")
		ar.Error(w, err, http.StatusBadRequest, "")
		return
	}

	extension := r.URL.Query().Get("ext")
	if len(extension) == 0 {
		err = fmt.Errorf("Empty extension")
		ar.Error(w, err, http.StatusBadRequest, "")
		return
	}

	filepath = strings.Join([]string{filename, extension}, ".")
	return
}
