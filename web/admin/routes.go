package admin

import (
	"github.com/gorilla/mux"
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

	// setup users routes
	users := mux.NewRouter().PathPrefix("/users").Subrouter()
	ar.router.PathPrefix("/users").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(users),
	))
	users.PathPrefix("/").HandlerFunc(ar.FetchUsers()).Methods("GET")
}
