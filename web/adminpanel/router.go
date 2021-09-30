package adminpanel

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
)

// Router is an admin panel router.
type Router struct {
	router *mux.Router
}

// NewRouter creates and initializes new admin panel router.
func NewRouter(options ...func(*Router) error) (model.Router, error) {
	apr := &Router{
		router: mux.NewRouter(),
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
