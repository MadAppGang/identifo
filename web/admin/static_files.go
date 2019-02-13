package admin

import (
	"html/template"
	"net/http"
	"path"
)

// StaticPages holds together all paths to static pages.
type StaticPages struct {
	AdminLogin string
}

// StaticFilesPath holds paths to static files.
type StaticFilesPath struct {
	StylesPath  string
	ScriptsPath string
	PagesPath   string
	ImagesPath  string
}

var defaultStaticPath = StaticFilesPath{
	StylesPath:  "../../static/css",
	ScriptsPath: "../../static/js",
	PagesPath:   "../../static",
	ImagesPath:  "../../static/img",
}

var defaultStaticPages = StaticPages{
	AdminLogin: "login.html", // TODO: change to "admin_login.html" when it is implemented.
}

// DefaultStaticPagesOptions sets default HTML pages.
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

// LoginHandler receives path to the admin login page template and serves it over http.
func (ar *Router) LoginHandler(pathComponents ...string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(pathComponents...))

	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		data := map[string]interface{}{
			"Prefix": ar.PathPrefix,
		} // TODO: fill when template is known.

		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
	}
}
