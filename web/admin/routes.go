package admin

import (
	"github.com/urfave/negroni"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	// do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	ar.router.Path("/{login:login\\/?}").Handler(negroni.New(
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")

	ar.router.Path("/{users:users\\/?}").Handler(ar.middleware.With(
		ar.Session(),
		negroni.WrapFunc(ar.FetchUsers()),
	)).Methods("GET")

	// setup users routes
	users := ar.router.PathPrefix("/users").Subrouter()
	ar.router.PathPrefix("/users").Handler(ar.middleware.With(
		ar.Session(),
		negroni.Wrap(users),
	))
}
