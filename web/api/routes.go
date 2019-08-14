package api

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// setup all routes
func (ar *Router) initRoutes() {
	if ar.router == nil {
		panic("Empty API router")
	}

	// All requests to the API router should contain appID.
	apiMiddlewares := ar.middleware.With(ar.DumpRequest(), ar.AppID())

	ar.router.HandleFunc(`/{ping:ping/?}`, ar.HandlePing()).Methods("GET")

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

	auth.Path(`/{tfa/enable:tfa/enable/?}`).Handler(negroni.New(
		ar.Token(TokenTypeAccess),
		negroni.Wrap(ar.EnableTFA()),
	)).Methods("PUT")
	auth.Path(`/{tfa/disable:tfa/disable/?}`).Handler(negroni.New(
		negroni.Wrap(ar.RequestDisabledTFA()),
	)).Methods("PUT")
	auth.Path(`/{tfa/finalize:tfa/finalize/?}`).Handler(negroni.New(
		ar.Token(TokenTypeAccess),
		negroni.Wrap(ar.FinalizeTFA()),
	)).Methods("POST")
	auth.Path(`/{tfa/reset:tfa/reset/?}`).Handler(negroni.New(
		ar.Token(TokenTypeAccess),
		negroni.Wrap(ar.RequestTFAReset()),
	)).Methods("PUT")

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	ar.router.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(TokenTypeAccess),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("").HandlerFunc(ar.IsLoggedIn()).Methods("GET")
	meRouter.Path("").HandlerFunc(ar.UpdateUser()).Methods("PUT")
	meRouter.Path(`/{logout:logout/?}`).HandlerFunc(ar.Logout()).Methods("POST")

	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	ar.router.PathPrefix("/.well-known").Handler(negroni.New(
		ar.DumpRequest(),
		negroni.Wrap(oidc),
	))

	oidc.Path(`/{openid-configuration:openid-configuration/?}`).HandlerFunc(ar.OIDCConfiguration()).Methods("GET")
	oidc.Path(`/{jwks.json:jwks.json/?}`).HandlerFunc(ar.OIDCJwks()).Methods("GET")
	oidc.Path(`/{apple-developer-domain-association.txt:apple-developer-domain-association.txt/?}`).HandlerFunc(ar.SupportSignInWithApple()).Methods("GET")
}
