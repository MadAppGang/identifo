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
	ar.handler.HandleFunc("/{ping:ping\\/?}", ar.HandlePing()).Methods("GET")
	ar.handler.HandleFunc("/password/{forgot:forgot\\/?}", ar.SendResetToken()).Methods("POST")

	//setup auth routes
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	ar.handler.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))
	auth.Path("/{login:login\\/?}").HandlerFunc(ar.LoginWithPassword()).Methods("POST")
	auth.Path("/{federated:federated\\/?}").HandlerFunc(ar.FederatedLogin()).Methods("POST")
	auth.Path("/{register:register\\/?}").HandlerFunc(ar.RegisterWithPassword()).Methods("POST")

	auth.Path("/{token:token\\/?}").Handler(negroni.New(
		ar.Token(TokenTypeRefresh),
		negroni.Wrap(ar.RefreshToken()),
	)).Methods("GET")

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	ar.handler.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(TokenTypeAccess),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("/{logout:logout\\/?}").HandlerFunc(ar.Logout()).Methods("POST")

	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	ar.handler.PathPrefix("/.well-known").Handler(negroni.New(
		ar.DumpRequest(),
		negroni.Wrap(oidc),
	))

	oidc.Path("/{openid-configuration:openid-configuration\\/?}").HandlerFunc(ar.OIDCConfiguration()).Methods("GET")
	oidc.Path("/{jwks.json:jwks.json\\/?}").HandlerFunc(ar.OIDCJwks()).Methods("GET")

}
