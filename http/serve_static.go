package http

import (
	"html/template"
	"net/http"
)

// StaticPages holds together all paths to a static pages
type StaticPages struct {
	Login                 string
	Registration          string
	ForgotPassword        string
	ForgotPasswordSuccess string
	ResetPassword         string
}

// StaticFiles holds paths to static files
type StaticFiles struct {
	StylesDirectory  string
	ScriptsDirectory string
}

// ServeTemplate receives path to a template and serves it over http
func (ar *apiRouter) ServeTemplate(path string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path)
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		errorMessage, err := GetFlash(w, r, ErrorMessageKey)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		data := map[string]interface{}{"Error": errorMessage}
		err = tmpl.Execute(w, data)
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

// ServeStaticFiles serves styles and scripts from provided directories
func ServeStaticFiles(sf StaticFiles) func(*apiRouter) error {
	return func(ar *apiRouter) error {
		return ar.serveStaticFiles(sf)
	}
}

// ServeDefaultStaticPages serves default HTML pages
func ServeDefaultStaticPages() func(*apiRouter) error {
	staticPages := StaticPages{
		Login:                 "./static/login.html",
		Registration:          "./static/registration.html",
		ForgotPassword:        "./static/forgot-password.html",
		ResetPassword:         "./static/reset-password.html",
		ForgotPasswordSuccess: "./static/forgot-password-success.html",
	}

	return ServeStaticPages(staticPages)
}

func (ar *apiRouter) serveStaticPages(sp StaticPages) error {
	//setup routes for templates
	ar.handler.HandleFunc("/{login:login\\/?}", ar.ServeTemplate(sp.Login)).Methods("GET")
	ar.handler.HandleFunc("/{register:register\\/?}", ar.ServeTemplate(sp.Registration)).Methods("GET")
	ar.handler.HandleFunc("/password/{forgot:forgot\\/?}", ar.ServeTemplate(sp.ForgotPassword)).Methods("GET")
	ar.handler.HandleFunc("/password/forgot/{success:success\\/?}", ar.ServeTemplate(sp.ForgotPasswordSuccess)).Methods("GET")

	return nil
}

// ServeDefaultStaticFiles serves styles and scripts from default path
func ServeDefaultStaticFiles() func(*apiRouter) error {
	staticFiles := StaticFiles{
		StylesDirectory:  "./static/css",
		ScriptsDirectory: "./static/js",
	}

	return ServeStaticFiles(staticFiles)
}

func (ar *apiRouter) serveStaticFiles(sf StaticFiles) error {
	stylesHandler := http.FileServer(http.Dir(sf.StylesDirectory))
	scriptsHandler := http.FileServer(http.Dir(sf.ScriptsDirectory))

	//setup routes for static files
	ar.handler.PathPrefix("/{css:css\\/?}").Handler(http.StripPrefix("/css/", stylesHandler))
	ar.handler.PathPrefix("/{js:js\\/?}").Handler(http.StripPrefix("/js/", scriptsHandler))

	return nil
}
