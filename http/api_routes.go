package http

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//setup all routes
func (ar *apiRouter) initRoutes() {
	//do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	//all API routes should have appID in it
	apiMiddlewares := ar.router.With(ar.DumpRequest(), ar.AppID())

	//setup root routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ping", ar.HandlePing()).Methods("GET")

	//setup auth routes
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	r.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))
	auth.Path("/login").HandlerFunc(ar.LoginWithPassword()).Methods("POST")

	ar.router.UseHandler(r)
}
