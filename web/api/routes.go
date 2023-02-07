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

	ar.router.HandleFunc("/ping", ar.HandlePing()).Methods("GET")

	apiMiddlewares := ar.buildAPIMiddleware()

	// federated oidc
	federatedOIDC := ar.router.PathPrefix("/auth/federated/oidc").Subrouter()
	ar.buildFederatedOIDCRoutes(federatedOIDC, apiMiddlewares)

	// auth
	auth := ar.buildAuthRoutes(apiMiddlewares)
	ar.router.PathPrefix("/auth").Handler(auth)

	// me
	me := ar.buildMeRoutes(apiMiddlewares)
	ar.router.PathPrefix("/me").Handler(me)

	// oidc config provider
	oidcMiddlewares := ar.buildOIDCMiddleware()
	oidcCfg := ar.buildOIDCRoutes(oidcMiddlewares)
	ar.router.PathPrefix("/.well-known").Handler(oidcCfg)
}

// buildAPIMiddleware creates middlewares that should execute for all requests
func (ar *Router) buildAPIMiddleware() *negroni.Negroni {
	handlers := []negroni.Handler{ar.ConfigCheck()}

	if ar.LoggerSettings.DumpRequest {
		handlers = append(handlers, ar.DumpRequest())
	}

	handlers = append(handlers, ar.AppID())
	return ar.middleware.With(handlers...)
}

func (ar *Router) buildFederatedOIDCRoutes(router *mux.Router, middlewares *negroni.Negroni) {
	router.Use(func(h http.Handler) http.Handler {
		return middlewares.With(negroni.Wrap(h))
	})

	router.Path("/login").Methods("POST").HandlerFunc(ar.OIDCLogin)
	router.Path("/login").Methods("GET").HandlerFunc(ar.OIDCLogin)

	// some OIDC providers do not allow to redirect to url with query params
	// so we have to use path argument to pass app id
	// it will not work with auth router since AppID middleware
	// will not be able to find app by id
	router.Path("/complete/{appId}").Methods("POST").HandlerFunc(ar.OIDCLoginComplete)
	router.Path("/complete/{appId}").Methods("GET").HandlerFunc(ar.OIDCLoginComplete)

	router.Path("/complete").HandlerFunc(ar.OIDCLoginComplete).Methods("POST")
	router.Path("/complete").HandlerFunc(ar.OIDCLoginComplete).Methods("GET")
}

func (ar *Router) buildAuthRoutes(middlewares *negroni.Negroni) *negroni.Negroni {
	// now build auth router for main API
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()

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

	return middlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	)
}

func (ar *Router) buildMeRoutes(middleware *negroni.Negroni) *negroni.Negroni {
	me := mux.NewRouter().PathPrefix("/me").Subrouter()
	me.Path("").HandlerFunc(ar.GetUser()).Methods("GET")
	me.Path("").HandlerFunc(ar.UpdateUser()).Methods("PUT")
	me.Path("/logout").HandlerFunc(ar.Logout()).Methods("POST")

	return middleware.With(
		ar.SignatureHandler(),
		ar.Token(model.TokenTypeAccess, nil),
		negroni.Wrap(me),
	)
}

func (ar *Router) buildOIDCRoutes(middleware *negroni.Negroni) *negroni.Negroni {
	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	oidc.Path("/openid-configuration").HandlerFunc(ar.OIDCConfiguration()).Methods("GET")
	oidc.Path("/jwks.json").HandlerFunc(ar.OIDCJwks()).Methods("GET")

	// apple native integration
	// TODO: Jack reimplement it completely
	// oidc.Path("/apple-developer-domain-association.txt").HandlerFunc(ar.ServeADDAFile()).Methods("GET")
	// oidc.Path("/apple-app-site-association").HandlerFunc(ar.ServeAASAFile()).Methods("GET")

	return ar.middleware.
		With(middleware.Handlers()...).
		With(negroni.Wrap(oidc))
}

func (ar *Router) buildOIDCMiddleware() *negroni.Negroni {
	wellKnownHandlers := make([]negroni.Handler, 0)
	wellKnownHandlers = append(wellKnownHandlers, ar.ConfigCheck())

	if ar.LoggerSettings.DumpRequest {
		wellKnownHandlers = append(wellKnownHandlers, ar.DumpRequest())
	}

	return ar.middleware.With(wellKnownHandlers...)
}
