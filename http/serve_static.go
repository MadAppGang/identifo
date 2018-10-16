package http

import (
	"html/template"
	"net/http"
)

// StaticPages holds together all static pages
type StaticPages struct {
	Login          string
	Registration   string
	ForgotPassword string
	ResetPassword  string
}

// ServeTemplate receives path to a template and serves it over http
func (ar *apiRouter) ServeTemplate(path string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path)
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		err := tmpl.Execute(w, nil)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
