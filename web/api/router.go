package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	r "github.com/madappgang/identifo/v2/web/router"
)

// Router is a router that handles all API requests.
type Router struct {
	r.LocalizedRouter
	server            model.Server
	cors              cors.Options
	router            *chi.Mux
	oidcConfiguration *OIDCConfiguration
	jwk               *JWK
	LoggerSettings    model.LoggerSettings
}

type RouterSettings struct {
	Server         model.Server
	Logger         *log.Logger
	LoggerSettings model.LoggerSettings
	Cors           cors.Options
	Locale         string
}

// NewRouter creates and inits new router.
func NewRouter(settings RouterSettings) (*Router, error) {
	l, err := l.NewPrinter(settings.Locale)
	if err != nil {
		return nil, err
	}

	ar := Router{
		server:         settings.Server,
		LoggerSettings: settings.LoggerSettings,
		cors:           settings.Cors,
	}

	ar.LP = l
	// setup logger to stdout.
	if settings.Logger == nil {
		ar.Logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		ar.Logger = settings.Logger
	}

	ar.initRoutes()
	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}
