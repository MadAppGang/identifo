package http

import (
	"html/template"
	"net/http"
)

// StaticPages holds together all paths to a static pages
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

// ServeStaticPages serves static provided pages
func ServeStaticPages(sp StaticPages) func(*apiRouter) error {
	return func(ar *apiRouter) error {
		return ar.serveStaticPages(sp)
	}
}

// ServeDefaultStaticPages serves default HTML pages
func ServeDefaultStaticPages() func(*apiRouter) error {
	staticPages := StaticPages{
		Login:          "../../static/login.html",
		Registration:   "../../static/registration.html",
		ForgotPassword: "../../static/forgot-password.html",
		ResetPassword:  "../../static/reset-password.html",
	}

	return ServeStaticPages(staticPages)
}

func (ar *apiRouter) serveStaticPages(sp StaticPages) error {
	ar.handler.HandleFunc("/login", ar.ServeTemplate(sp.Login)).Methods("GET")
	ar.handler.HandleFunc("/register", ar.ServeTemplate(sp.Registration)).Methods("GET")
	ar.handler.HandleFunc("/password/forgot", ar.ServeTemplate(sp.ForgotPassword)).Methods("GET")
	ar.handler.HandleFunc("/password/reset", ar.ServeTemplate(sp.ResetPassword)).Methods("GET")
	return nil
}
