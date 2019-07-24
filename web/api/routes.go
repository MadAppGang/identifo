package api

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// setup all routes
func (ar *Router) initRoutes() {
	// do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	// all API routes should have appID in it
	apiMiddlewares := ar.middleware.With(ar.DumpRequest(), ar.AppID())

	// setup root routes
	ar.router.HandleFunc(`/{ping:ping/?}`, ar.HandlePing()).Methods("GET")

	// setup auth routes
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	ar.router.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))
	auth.Path(`/{login:login/?}`).HandlerFunc(ar.LoginWithPassword()).Methods("POST")
	auth.Path(`/{request_phone_code:request_phone_code/?}`).HandlerFunc(ar.RequestVerificationCode()).Methods("POST")
	auth.Path(`/{phone_login:phone_login/?}`).HandlerFunc(ar.PhoneLogin()).Methods("POST")
	auth.Path(`/{federated:federated/?}`).HandlerFunc(ar.FederatedLogin()).Methods("POST")
	auth.Path(`/{register:register/?}`).HandlerFunc(ar.RegisterWithPassword()).Methods("POST")
	auth.Path(`/{reset_password:reset_password/?}`).HandlerFunc(ar.RequestResetPassword()).Methods("POST")

	auth.Path(`/{token:token/?}`).Handler(negroni.New(
		ar.Token(TokenTypeRefresh),
		negroni.Wrap(ar.RefreshTokens()),
	)).Methods("POST")
	auth.Path(`/{invite:invite/?}`).Handler(negroni.New(
		ar.Token(TokenTypeAccess),
		negroni.Wrap(ar.RequestInviteLink()),
	)).Methods("POST")

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	ar.router.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(TokenTypeAccess),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("").HandlerFunc(ar.UpdateUser()).Methods("PUT")
	meRouter.Path(`/{logout:logout/?}`).HandlerFunc(ar.Logout()).Methods("POST")

	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	ar.router.PathPrefix("/.well-known").Handler(negroni.New(
		ar.DumpRequest(),
		negroni.Wrap(oidc),
	))

	oidc.Path(`/{openid-configuration:openid-configuration/?}`).HandlerFunc(ar.OIDCConfiguration()).Methods("GET")
	oidc.Path(`/{jwks.json:jwks.json/?}`).HandlerFunc(ar.OIDCJwks()).Methods("GET")
}
