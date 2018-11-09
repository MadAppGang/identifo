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
	ar.handler.HandleFunc("/ping", ar.HandlePing()).Methods("GET")

	//setup auth routes
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	ar.handler.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))
	auth.Path("/login").HandlerFunc(ar.LoginWithPassword()).Methods("POST")
	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods("POST")
	auth.Path("/register").HandlerFunc(ar.RegisterWithPassword()).Methods("POST")

	auth.Path("/token").Handler(negroni.New(
		ar.Token(TokenTypeRefresh),
		negroni.Wrap(ar.RefreshToken()),
	)).Methods("GET")

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	ar.handler.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(TokenTypeAccess),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("/logout").HandlerFunc(ar.Logout()).Methods("POST")
}
