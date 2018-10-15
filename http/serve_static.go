package http

import (
	"html/template"
	"net/http"
)

func (ar *apiRouter) ServeTemplate(path string) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(path))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
