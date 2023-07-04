package api

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	r "github.com/madappgang/identifo/v2/web/router"
	"github.com/rs/cors"
)

// Router is a router that handles all API requests.
type Router struct {
	r.LocalizedRouter
	server               model.Server
	cors                 *cors.Cors
	router               *mux.Router
	tfaType              model.TFAType
	tfaResendTimeout     int
	oidcConfiguration    *OIDCConfiguration
	jwk                  *JWK
	Authorizer           *authorization.Authorizer
	Host                 *url.URL
	SupportedLoginWays   model.LoginWith
	tokenPayloadServices map[string]model.TokenPayloadProvider
	LoggerSettings       model.LoggerSettings
}

type RouterSettings struct {
	Server           model.Server
	Logger           *log.Logger
	LoggerSettings   model.LoggerSettings
	Authorizer       *authorization.Authorizer
	Host             *url.URL
	TFAType          model.TFAType
	TFAResendTimeout int
	LoginWith        model.LoginWith
	Cors             *cors.Cors
	Locale           string
}

// NewRouter creates and inits new router.
func NewRouter(settings RouterSettings) (*Router, error) {
	l, err := l.NewPrinter(settings.Locale)
	if err != nil {
		return nil, err
	}

	ar := Router{
		server:             settings.Server,
		router:             mux.NewRouter(),
		Authorizer:         settings.Authorizer,
		LoggerSettings:     settings.LoggerSettings,
		Host:               settings.Host,
		tfaType:            settings.TFAType,
		tfaResendTimeout:   settings.TFAResendTimeout,
		SupportedLoginWays: settings.LoginWith,
		cors:               settings.Cors,
	}

	ar.LP = l
	// setup logger to stdout.
	if settings.Logger == nil {
		ar.Logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		ar.Logger = settings.Logger
	}

	ar.tokenPayloadServices = make(map[string]model.TokenPayloadProvider)

	ar.initRoutes()

	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}
