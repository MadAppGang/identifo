package html

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router handles incoming http connections.
type Router struct {
	Server     model.Server
	Middleware *negroni.Negroni
	Logger     *log.Logger
	Router     *mux.Router
	Authorizer *authorization.Authorizer
	PathPrefix string
	Host       string
	cors       *cors.Cors
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		PathPrefixOptions("/web"),
	}
}

// PathPrefixOptions set path prefix options.
func PathPrefixOptions(prefix string) func(r *Router) error {
	return func(r *Router) error {
		r.PathPrefix = prefix
		return nil
	}
}

// HostOption sets hostname.
func HostOption(host string) func(r *Router) error {
	return func(r *Router) error {
		r.Host = host
		return nil
	}
}

// CorsOption sets cors option.
func CorsOption(corsOptions model.CorsOptions) func(*Router) error {
	return func(r *Router) error {
		if corsOptions.HTML != nil {
			r.cors = cors.New(*corsOptions.HTML)
		}
		return nil
	}
}

// NewRouter creates and initializes new router.
func NewRouter(server model.Server, logger *log.Logger, authorizer *authorization.Authorizer, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		Middleware: negroni.Classic(),
		Router:     mux.NewRouter(),
		Server:     server,
		Authorizer: authorizer,
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	// Setup logger to stdout.
	if logger == nil {
		ar.Logger = log.New(os.Stdout, "HTML_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if ar.cors != nil {
		ar.Middleware.Use(ar.cors)
	}
	ar.initRoutes()
	ar.Middleware.UseHandler(ar.Router)

	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.Router.ServeHTTP(w, r)
}
