package html

import (
	"html/template"
	"net/http"
	"path"

	"github.com/madappgang/identifo/model"
)

// StaticPages holds together all paths to a static pages
type StaticPages struct {
	Login                 string
	Registration          string
	ForgotPassword        string
	ForgotPasswordSuccess string
	ResetPassword         string
	TokenError            string
	ResetSuccess          string
	Misconfiguration      string
}

//EmailTemplates store email templates
type EmailTemplates struct {
	Welcome       string
	ResetPassword string
	VerifyEmail   string
}

// StaticFilesPath holds paths to static files
type StaticFilesPath struct {
	StylesPath     string
	ScriptsPath    string
	PagesPath      string
	ImagesPath     string
	EmailTemplates string
}

var defaultStaticPath = StaticFilesPath{
	StylesPath:  "./static/css",
	ScriptsPath: "./static/js",
	PagesPath:   "./static",
	ImagesPath:  "./static/img",
}

var defaultStaticPages = StaticPages{
	Login:                 "login.html",
	Registration:          "registration.html",
	ForgotPassword:        "forgot-password.html",
	ResetPassword:         "reset-password.html",
	ForgotPasswordSuccess: "forgot-password-success.html",
	TokenError:            "token-error.html",
	ResetSuccess:          "reset-success.html",
	Misconfiguration:      "misconfiguration.html",
}

// DefaultStaticPagesOptions set default HTML pages
func DefaultStaticPagesOptions() func(r *Router) error {
	return func(r *Router) error {
		r.StaticPages = defaultStaticPages
		return nil
	}
}

// DefaultStaticPathOptions set default static files locations
func DefaultStaticPathOptions() func(r *Router) error {
	return func(r *Router) error {
		r.StaticFilesPath = defaultStaticPath
		return nil
	}
}

// StaticPathOptions set  static files locations
func StaticPathOptions(path StaticFilesPath) func(r *Router) error {
	return func(r *Router) error {
		r.StaticFilesPath = path
		return nil
	}
}

// HTMLFileHandler receives path to a template and serves it over http
func (ar *Router) HTMLFileHandler(pathComponents ...string) http.HandlerFunc {

	tmpl, err := template.ParseFiles(path.Join(pathComponents...))
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
		err = tmpl.Execute(w, data)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}

}

//ResetPasswordHandler handles reset password request
func (ar *Router) ResetPasswordHandler(pathComponents ...string) http.HandlerFunc {

	tmpl, err := template.ParseFiles(path.Join(pathComponents...))

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

		token := r.Context().Value(model.TokenRawContextKey)
		data := map[string]interface{}{
			"Error":  errorMessage,
			"Token":  token,
			"Prefix": ar.PathPrefix,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}

// LoginHandler handles login page request.
func (ar *Router) LoginHandler(pathComponents ...string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(pathComponents...))

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

		scopes := r.URL.Query().Get("scopes")
		app := appFromContext(r.Context())

		data := map[string]interface{}{
			"Error":  errorMessage,
			"Prefix": ar.PathPrefix,
			"Scopes": scopes,
			"AppId":  app.ID(),
		}

		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
