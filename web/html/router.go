package html

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router handles incoming http connections.
type Router struct {
	Middleware         *negroni.Negroni
	Logger             *log.Logger
	Router             *mux.Router
	AppStorage         model.AppStorage
	UserStorage        model.UserStorage
	TokenStorage       model.TokenStorage
	TokenBlacklist     model.TokenBlacklist
	TokenService       jwtService.TokenService
	SMSService         model.SMSService
	EmailService       model.EmailService
	staticFilesStorage model.StaticFilesStorage
	Authorizer         *authorization.Authorizer
	PathPrefix         string
	Host               string
	cors               *cors.Cors
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
func CorsOption(corsOptions *model.CorsOptions) func(*Router) error {
	return func(r *Router) error {
		if corsOptions != nil && corsOptions.HTML != nil {
			r.cors = cors.New(*corsOptions.HTML)
		}
		return nil
	}
}

// NewRouter creates and initializes new router.
func NewRouter(logger *log.Logger, as model.AppStorage, us model.UserStorage, sfs model.StaticFilesStorage, ts model.TokenStorage, tb model.TokenBlacklist, tServ jwtService.TokenService, smsServ model.SMSService, emailServ model.EmailService, authorizer *authorization.Authorizer, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		Middleware:         negroni.Classic(),
		Router:             mux.NewRouter(),
		AppStorage:         as,
		UserStorage:        us,
		TokenStorage:       ts,
		TokenBlacklist:     tb,
		TokenService:       tServ,
		SMSService:         smsServ,
		EmailService:       emailServ,
		staticFilesStorage: sfs,
		Authorizer:         authorizer,
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
	ar.Middleware.UseHandler(ar.Router)

	ar.initRoutes()
	return &ar, nil
}

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, err error, code int, userInfo string) {
	// Log error.
	ar.Logger.Printf("http error: %s (code=%d)\n", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = model.ErrorInternal
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseString := `
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Home Network</title>
	</head>
	<body>
	<h2>Error</h2></br>
	<h3>
	` +
		fmt.Sprintf("Error: %s, code: %d, userInfo: %s", err.Error(), code, userInfo) +
		`
	</h3>
	</body>
	</html>
	`
	w.WriteHeader(code)
	if _, wrErr := io.WriteString(w, responseString); wrErr != nil {
		ar.Logger.Println("Error writing response string:", wrErr)
	}
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.Router.ServeHTTP(w, r)
}
