package admin

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

// Router is a router that handles admin requests.
type Router struct {
	r.LocalizedRouter
	server       model.Server
	cors         cors.Options
	router       *chi.Mux
	forceRestart chan<- bool
	originUpdate func() error
}

type RouterSettings struct {
	Server       model.Server
	Logger       *log.Logger
	Cors         cors.Options
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
		forceRestart: settings.Restart,
		cors:         settings.Cors,
		originUpdate: settings.OriginUpdate,
	}

	if settings.Logger == nil {
		ar.Logger = log.New(os.Stdout, "ADMIN_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	ar.LP = l
	ar.initRoutes()
	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.router.ServeHTTP(w, r)
}
