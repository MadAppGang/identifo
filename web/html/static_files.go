package html

import (
	"html/template"
	"net/http"
	"path"
)

// StaticPages holds together all paths to static pages.
type StaticPages struct {
	Login                 string
	Registration          string
	ForgotPassword        string
	ForgotPasswordSuccess string
	ResetPassword         string
	ResetPasswordSuccess  string
	DisableTFA            string
	DisableTFASuccess     string
	ResetTFA              string
	ResetTFASuccess       string
	TokenError            string
	WebMessage            string
	Misconfiguration      string
}

// EmailTemplates stores email templates.
type EmailTemplates struct {
	Welcome       string
	ResetPassword string
	VerifyEmail   string
}

// StaticFilesPath holds paths to static files.
type StaticFilesPath struct {
	StylesPath     string
	ScriptsPath    string
	PagesPath      string
	ImagesPath     string
	FontsPath      string
	EmailTemplates string
}

var defaultStaticPath = StaticFilesPath{
	StylesPath:  "./static/css",
	ScriptsPath: "./static/js",
	ImagesPath:  "./static/img",
	FontsPath:   "./static/fonts",
	PagesPath:   "./static",
}

var defaultStaticPages = StaticPages{
	Login:                 "login.html",
	Registration:          "registration.html",
	ForgotPassword:        "forgot-password.html",
	ResetPassword:         "reset-password.html",
	ResetPasswordSuccess:  "reset-password-success.html",
	DisableTFA:            "disable-tfa.html",
	DisableTFASuccess:     "disable-tfa-success.html",
	ResetTFA:              "reset-tfa.html",
	ResetTFASuccess:       "reset-tfa-success.html",
	ForgotPasswordSuccess: "forgot-password-success.html",
	TokenError:            "token-error.html",
	WebMessage:            "web-message.html",
	Misconfiguration:      "misconfiguration.html",
}

// DefaultStaticPagesOptions sets default HTML pages path.
func DefaultStaticPagesOptions() func(r *Router) error {
	return func(r *Router) error {
		r.StaticPages = defaultStaticPages
		return nil
	}
}

// DefaultStaticPathOptions sets default static files locations.
func DefaultStaticPathOptions() func(r *Router) error {
	return func(r *Router) error {
		r.StaticFilesPath = defaultStaticPath
		return nil
	}
}

// StaticPathOptions sets static files locations.
func StaticPathOptions(path StaticFilesPath) func(r *Router) error {
	return func(r *Router) error {
		r.StaticFilesPath = path
		return nil
	}
}

// HTMLFileHandler receives path to a template and serves it over HTTP.
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
		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
