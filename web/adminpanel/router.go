package adminpanel

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
)

// Router is an admin panel router.
type Router struct {
	router    *mux.Router
	buildPath string
}

// NewRouter creates and initializes new admin panel router.
func NewRouter(buildPath string, options ...func(*Router) error) (model.Router, error) {
	apr := &Router{
		router:    mux.NewRouter(),
		buildPath: buildPath,
	}

	for _, option := range options {
		if err := option(apr); err != nil {
			return nil, err
		}
	}

	apr.initRoutes()
	return apr, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.router.ServeHTTP(w, r)
}

// BuildPathOption sets admin panel build path.
func BuildPathOption(buildpath string) func(*Router) error {
	return func(r *Router) error {
		r.buildPath = buildpath
		return nil
	}
}
