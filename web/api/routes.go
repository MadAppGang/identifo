package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/model"
	"github.com/urfave/negroni"
)

// setup all routes
func (ar *Router) initRoutes() {
	if ar.router == nil {
		panic("Empty API router")
	}

	// All requests to the API router should contain appID.
	handlers := make([]negroni.Handler, 0)
	handlers = append(handlers, ar.ConfigCheck())

	if ar.LoggerSettings.DumpRequest {
		handlers = append(handlers, ar.DumpRequest())
	}

	handlers = append(handlers, ar.AppID())

	apiMiddlewares := ar.middleware.With(handlers...)

	ar.router.HandleFunc("/ping", ar.HandlePing()).Methods("GET")

	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	ar.router.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))

	auth.Path("/login").HandlerFunc(ar.LoginWithPassword()).Methods("POST")
	auth.Path("/request_phone_code").HandlerFunc(ar.RequestVerificationCode()).Methods("POST")
	auth.Path("/phone_login").HandlerFunc(ar.PhoneLogin()).Methods("POST")
	auth.Path("/register").HandlerFunc(ar.RegisterWithPassword()).Methods("POST")
	auth.Path("/request_reset_password").HandlerFunc(ar.RequestResetPassword()).Methods("POST")
	auth.Path("/reset_password").Handler(negroni.New(
		ar.Token(model.TokenTypeReset, nil),
		negroni.Wrap(ar.ResetPassword()),
	)).Methods("POST")

	auth.Path("/app_settings").HandlerFunc(ar.GetAppSettings()).Methods("GET")

	auth.Path("/token").Handler(negroni.New(
		ar.Token(model.TokenTypeRefresh, nil),
		negroni.Wrap(ar.RefreshTokens()),
	)).Methods("POST")
	auth.Path("/invite").Handler(negroni.New(
		ar.Token(model.TokenTypeAccess, nil),
		negroni.Wrap(ar.RequestInviteLink()),
	)).Methods("POST")

	auth.Path("/tfa/enable").Handler(negroni.New(
		ar.Token(model.TokenTypeAccess, nil),
		negroni.Wrap(ar.EnableTFA()),
	)).Methods("PUT")
	auth.Path("/tfa/disable").Handler(negroni.New(
		negroni.Wrap(ar.RequestDisabledTFA()),
	)).Methods("PUT")
	auth.Path("/tfa/login").Handler(negroni.New(
		ar.Token(model.TokenTypeAccess, []string{model.TokenTypeTFAPreauth}),
		negroni.Wrap(ar.FinalizeTFA()),
	)).Methods("POST")
	auth.Path("/tfa/resend").Handler(negroni.New(
		ar.Token(model.TokenTypeAccess, []string{model.TokenTypeTFAPreauth}),
		negroni.Wrap(ar.ResendTFA()),
	)).Methods("POST")
	auth.Path("/tfa/reset").Handler(negroni.New(
		ar.Token(model.TokenTypeAccess, nil),
		negroni.Wrap(ar.RequestTFAReset()),
	)).Methods("PUT")

	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods("POST")
	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods("GET")

	auth.Path("/federated/complete").HandlerFunc(ar.FederatedLoginComplete()).Methods("POST")
	auth.Path("/federated/complete").HandlerFunc(ar.FederatedLoginComplete()).Methods("GET")

	auth.Path("/federated/oidc").HandlerFunc(ar.OIDCLogin).Methods("POST")
	auth.Path("/federated/oidc").HandlerFunc(ar.OIDCLogin).Methods("GET")

	auth.Path("/federated/oidc/complete").HandlerFunc(ar.OIDCLoginComplete).Methods("POST")
	auth.Path("/federated/oidc/complete").HandlerFunc(ar.OIDCLoginComplete).Methods("GET")

	fr := auth.Path("/federated/oidc/complete/{appId}").
		Subrouter()

	fr.Use(func(h http.Handler) http.Handler {
		m := ar.AppID()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m(w, r, h.ServeHTTP)
		})
	})

	fr.Methods("POST").HandlerFunc(ar.OIDCLoginComplete)
	fr.Methods("GET").HandlerFunc(ar.OIDCLoginComplete)

	meRouter := mux.NewRouter().PathPrefix("/me").Subrouter()
	ar.router.PathPrefix("/me").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		ar.Token(model.TokenTypeAccess, nil),
		negroni.Wrap(meRouter),
	))
	meRouter.Path("").HandlerFunc(ar.GetUser()).Methods("GET")
	meRouter.Path("").HandlerFunc(ar.UpdateUser()).Methods("PUT")
	meRouter.Path("/logout").HandlerFunc(ar.Logout()).Methods("POST")

	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	wellKnownHandlers := make([]negroni.Handler, 0)
	wellKnownHandlers = append(wellKnownHandlers, ar.ConfigCheck())

	if ar.LoggerSettings.DumpRequest {
		wellKnownHandlers = append(wellKnownHandlers, ar.DumpRequest())
	}

	wellKnownHandlers = append(wellKnownHandlers, negroni.Wrap(oidc))

	ar.router.PathPrefix("/.well-known").Handler(ar.middleware.With(wellKnownHandlers...))

	oidc.Path("/openid-configuration").HandlerFunc(ar.OIDCConfiguration()).Methods("GET")
	oidc.Path("/jwks.json").HandlerFunc(ar.OIDCJwks()).Methods("GET")

	// apple native integration
	// TODO: Jack reimplement it completely
	// oidc.Path("/apple-developer-domain-association.txt").HandlerFunc(ar.ServeADDAFile()).Methods("GET")
	// oidc.Path("/apple-app-site-association").HandlerFunc(ar.ServeAASAFile()).Methods("GET")
}
