package html

import (
	"net/http"
	"path"

	"github.com/madappgang/identifo/model"
)

// HTMLFileHandler receives path to a template and serves it over HTTP.
func (ar *Router) HTMLFileHandler(templateName model.TemplateName) http.HandlerFunc {
	tmpl, err := ar.staticFilesStorage.ParseTemplate(templateName)
	if err != nil {
		ar.Logger.Fatalf("Cannot parse %v template. %s\n", templateName, err)
	}
	prefix := path.Clean(ar.PathPrefix)

	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		data := map[string]interface{}{
			"Error":  errorMessage,
			"Prefix": prefix,
		}
		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
