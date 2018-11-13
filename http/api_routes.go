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
	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods("POST")
	auth.Path("/register").HandlerFunc(ar.RegisterWithPassword()).Methods("POST")

	auth.Path("/token").Handler(negroni.New(
		ar.Token(TokenTypeRefresh),
		negroni.Wrap(ar.RefreshToken()),
	)).Methods("GET")

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	r.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(TokenTypeAccess),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("/logout").HandlerFunc(ar.Logout()).Methods("POST")

	wellKnownRouter := r.PathPrefix("/.well-known").Subrouter()
	wellKnownRouter.HandleFunc(("/openid-configuration"), ar.Configuration()).Methods("GET")
	wellKnownRouter.HandleFunc("/jwks", ar.ServeJWKS()).Methods("GET")

	ar.router.UseHandler(r)
}
