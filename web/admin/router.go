package admin

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	r "github.com/madappgang/identifo/v2/web/router"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router is a router that handles admin requests.
type Router struct {
	r.LocalizedRouter
	server       model.Server
	middleware   *negroni.Negroni
	cors         *cors.Cors
	router       *mux.Router
	forceRestart chan<- bool
	originUpdate func() error
}

type RouterSettings struct {
	Server       model.Server
	Logger       *log.Logger
	Cors         *cors.Cors
	Restart      chan<- bool
	OriginUpdate func() error
	Locale       string
}

// NewRouter creates and initializes new admin router.
func NewRouter(settings RouterSettings) (model.Router, error) {
	l, err := l.NewPrinter(settings.Locale)
	if err != nil {
		return nil, err
	}

	ar := Router{
		server:       settings.Server,
		middleware:   negroni.New(middleware.NewNegroniLogger("ADMIN_API"), negroni.NewRecovery()),
		router:       mux.NewRouter(),
		forceRestart: settings.Restart,
		cors:         settings.Cors,
		originUpdate: settings.OriginUpdate,
	}

	if settings.Logger == nil {
		ar.Logger = log.New(os.Stdout, "ADMIN_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	ar.LP = l

	ar.middleware.Use(ar.RemoveTrailingSlash())

	if ar.cors == nil {
		ar.cors = cors.New(model.DefaultCors)
	}
	ar.middleware.Use(ar.cors)

	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.middleware.ServeHTTP(w, r)
}
