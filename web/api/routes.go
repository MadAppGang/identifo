package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// setup all routes
func (ar *Router) initRoutes(
	loggerSettings model.LoggerSettings,
) {
	if ar.router == nil {
		panic("Empty API router")
	}

	pingHandler := negroni.New(
		negroni.NewRecovery(),
		negroni.WrapFunc(ar.HandlePing),
	)
	ar.router.Handle("/ping", pingHandler).Methods(http.MethodGet)

	baseMiddleware := buildBaseMiddleware(
		loggerSettings.DumpRequest,
		loggerSettings.Format,
		loggerSettings.API,
		loggerSettings.LogSensitiveData,
		ar.cors,
	)
	apiMiddlewares := ar.buildAPIMiddleware(baseMiddleware)

	// federated oidc
	federatedOIDC := ar.router.PathPrefix("/auth/federated/oidc").Subrouter()
	ar.buildFederatedOIDCRoutes(federatedOIDC, apiMiddlewares, false)

	federatedOIDCV2 := ar.router.PathPrefix("/v2/auth/federated/oidc").Subrouter()
	ar.buildFederatedOIDCRoutes(federatedOIDCV2, apiMiddlewares, true)

	// auth
	auth := ar.buildAuthRoutes(apiMiddlewares)
	ar.router.PathPrefix("/auth").Handler(auth)

	// me
	me := ar.buildMeRoutes(apiMiddlewares)
	ar.router.PathPrefix("/me").Handler(me)

	// oidc config provider
	oidcMiddlewares := ar.buildOIDCMiddleware(baseMiddleware)
	oidcCfg := ar.buildOIDCRoutes(oidcMiddlewares)
	ar.router.PathPrefix("/.well-known").Handler(oidcCfg)
}

func buildBaseMiddleware(
	dumpRequest bool,
	format string,
	logParams model.LoggerParams,
	logSensitiveData bool,
	cors *cors.Cors,
) *negroni.Negroni {
	exclude := []string{}

	if !logSensitiveData {
		exclude = []string{
			"auth/login",
			"auth/request_phone_code",
			"auth/phone_login",
			"auth/register",
			"auth/token",
			"auth/request_reset_password",
			"auth/reset_password",
			"me/logout",
			"me/impersonate_as",
		}
	}

	lm := middleware.NegroniHTTPLogger(
		logging.ComponentAPI,
		format,
		logParams,
		model.HTTPLogDetailing(dumpRequest, logParams.HTTPDetailing),
		!logSensitiveData,
		exclude...)

	result := negroni.New(
		lm,
		negroni.NewRecovery(),
		middleware.RemoveTrailingSlash(),
	)

	if cors != nil {
		result.Use(cors)
	}

	return result

}

// buildAPIMiddleware creates middlewares that should execute for all requests
func (ar *Router) buildAPIMiddleware(base *negroni.Negroni) *negroni.Negroni {
	handlers := []negroni.Handler{ar.ConfigCheck()}
	handlers = append(handlers, ar.AppID())
	return with(base, handlers...)
}

func (ar *Router) buildFederatedOIDCRoutes(router *mux.Router, middlewares *negroni.Negroni, stateManagedByClient bool) {
	router.Use(func(h http.Handler) http.Handler {
		return with(middlewares, negroni.Wrap(h))
	})

	router.Path("/login").Methods(http.MethodPost).HandlerFunc(ar.OIDCLogin(stateManagedByClient))
	router.Path("/login").Methods(http.MethodGet).HandlerFunc(ar.OIDCLogin(stateManagedByClient))

	// some OIDC providers do not allow to redirect to url with query params
	// so we have to use path argument to pass app id
	// it will not work with auth router since AppID middleware
	// will not be able to find app by id
	router.Path("/complete/{appId}").Methods(http.MethodPost).HandlerFunc(ar.OIDCLoginComplete(!stateManagedByClient))
	router.Path("/complete/{appId}").Methods(http.MethodGet).HandlerFunc(ar.OIDCLoginComplete(!stateManagedByClient))

	router.Path("/complete").HandlerFunc(ar.OIDCLoginComplete(!stateManagedByClient)).Methods(http.MethodPost)
	router.Path("/complete").HandlerFunc(ar.OIDCLoginComplete(!stateManagedByClient)).Methods(http.MethodGet)
}

func (ar *Router) buildAuthRoutes(middlewares *negroni.Negroni) http.Handler {
	// now build auth router for main API
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()

	auth.Path("/login").HandlerFunc(ar.LoginWithPassword()).Methods(http.MethodPost)
	auth.Path("/request_phone_code").HandlerFunc(ar.RequestVerificationCode()).Methods(http.MethodPost)
	auth.Path("/phone_login").HandlerFunc(ar.PhoneLogin()).Methods(http.MethodPost)
	auth.Path("/register").HandlerFunc(ar.RegisterWithPassword()).Methods(http.MethodPost)
	auth.Path("/request_reset_password").HandlerFunc(ar.RequestResetPassword()).Methods(http.MethodPost)
	auth.Path("/reset_password").Handler(
		ar.Token(model.TokenTypeReset, nil)(ar.ResetPassword()),
	).Methods(http.MethodPost)

	auth.Path("/app_settings").HandlerFunc(ar.GetAppSettings()).Methods(http.MethodGet)

	auth.Path("/token").Handler(
		ar.Token(model.TokenTypeRefresh, nil)(ar.RefreshTokens()),
	).Methods(http.MethodPost)
	auth.Path("/invite").Handler(
		ar.Token(model.TokenTypeAccess, nil)(ar.RequestInviteLink()),
	).Methods(http.MethodPost)
	auth.Path("/impersonate_token").Handler(
		ar.GetImpersonateToken()).Methods(http.MethodPost)

	auth.Path("/tfa/enable").Handler(
		ar.Token(model.TokenTypeAccess, nil)(ar.EnableTFA()),
	).Methods(http.MethodPut)
	auth.Path("/tfa/disable").Handler(
		ar.RequestDisabledTFA(),
	).Methods(http.MethodPut)
	auth.Path("/tfa/login").Handler(
		ar.Token(model.TokenTypeAccess, []string{model.TokenTypeTFAPreauth})(ar.FinalizeTFA()),
	).Methods(http.MethodPost)
	auth.Path("/tfa/resend").Handler(
		ar.Token(model.TokenTypeAccess, []string{model.TokenTypeTFAPreauth})(ar.ResendTFA()),
	).Methods(http.MethodPost)
	auth.Path("/tfa/reset").Handler(
		ar.Token(model.TokenTypeAccess, nil)(ar.RequestTFAReset()),
	).Methods(http.MethodPut)

	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods(http.MethodPost)
	auth.Path("/federated").HandlerFunc(ar.FederatedLogin()).Methods(http.MethodGet)

	auth.Path("/federated/complete").HandlerFunc(ar.FederatedLoginComplete()).Methods(http.MethodPost)
	auth.Path("/federated/complete").HandlerFunc(ar.FederatedLoginComplete()).Methods(http.MethodGet)

	return with(middlewares,
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	)
}

func (ar *Router) buildMeRoutes(middleware *negroni.Negroni) http.Handler {
	me := mux.NewRouter().PathPrefix("/me").Subrouter()

	me.Path("").HandlerFunc(ar.GetUser()).Methods(http.MethodGet)
	me.Path("").HandlerFunc(ar.UpdateUser()).Methods(http.MethodPut)
	me.Path("/logout").HandlerFunc(ar.Logout()).Methods(http.MethodPost)
	me.Path("/impersonate_as").HandlerFunc(ar.ImpersonateAs()).Methods(http.MethodPost)

	return with(middleware,
		ar.SignatureHandler(),
		negroni.Wrap(ar.Token(model.TokenTypeAccess, nil)(me)),
	)
}

func (ar *Router) buildOIDCRoutes(middleware *negroni.Negroni) http.Handler {
	oidc := mux.NewRouter().PathPrefix("/.well-known").Subrouter()

	oidc.Path("/openid-configuration").HandlerFunc(ar.OIDCConfiguration()).Methods(http.MethodGet)
	oidc.Path("/jwks.json").HandlerFunc(ar.OIDCJwks()).Methods(http.MethodGet)

	// apple native integration
	// TODO: Jack reimplement it completely
	// oidc.Path("/apple-developer-domain-association.txt").HandlerFunc(ar.ServeADDAFile()).Methods(http.MethodGet)
	// oidc.Path("/apple-app-site-association").HandlerFunc(ar.ServeAASAFile()).Methods(http.MethodGet)

	return with(middleware, negroni.Wrap(oidc))
}

func (ar *Router) buildOIDCMiddleware(base *negroni.Negroni) *negroni.Negroni {
	wellKnownHandlers := []negroni.Handler{ar.ConfigCheck()}

	return with(base, wellKnownHandlers...)
}

func with(n *negroni.Negroni, handlers ...negroni.Handler) *negroni.Negroni {
	existing := n.Handlers()
	h := []negroni.Handler{}
	h = append(h, existing...)
	h = append(h, handlers...)
	return negroni.New(h...)
}
