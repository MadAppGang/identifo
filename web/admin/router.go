package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router is a router that handles admin requests.
type Router struct {
	server            model.Server
	middleware        *negroni.Negroni
	cors              *cors.Cors
	logger            *log.Logger
	router            *mux.Router
	RedirectURL       string
	PathPrefix        string
	LoginWebAppPrefix string
	Host              string
	forceRestart      chan<- bool
	originUpdate      func() error
}

type RouterSettings struct {
	Server            model.Server
	Logger            *log.Logger
	Host              string
	Prefix            string
	LoginWebAppPrefix string // this used to create invite link and reset password link for user
	Cors              *cors.Cors
	Restart           chan<- bool
	OriginUpdate      func() error
}

// NewRouter creates and initializes new admin router.
func NewRouter(settings RouterSettings) (model.Router, error) {
	ar := Router{
		server:            settings.Server,
		middleware:        negroni.New(middleware.NewNegroniLogger("ADMIN_API"), negroni.NewRecovery()),
		router:            mux.NewRouter(),
		Host:              settings.Host,
		PathPrefix:        settings.Prefix,
		forceRestart:      settings.Restart,
		cors:              settings.Cors,
		RedirectURL:       "/login",
		LoginWebAppPrefix: settings.LoginWebAppPrefix,
	}

	if settings.Logger == nil {
		ar.logger = log.New(os.Stdout, "ADMIN_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.middleware.Use(ar.RemoveTrailingSlash())

	if ar.cors == nil {
		ar.cors = cors.New(model.DefaultCors)
	}
	ar.middleware.Use(ar.cors)

	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
}

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, err error, code int, userInfo string) {
	// errorResponse is a generic response for sending errors.
	type errorResponse struct {
		Error string `json:"error,omitempty"`
		Info  string `json:"info,omitempty"`
		Code  int    `json:"code,omitempty"`
	}

	// Log error.
	ar.logger.Printf("admin error: %v (code=%d)", err, code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encodeErr := json.NewEncoder(w).Encode(&errorResponse{
		Error: err.Error(),
		Info:  userInfo,
		Code:  code,
	})
	if encodeErr != nil {
		ar.logger.Printf("error writing http response: %s", err)
	}
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.middleware.ServeHTTP(w, r)
}

// ServeJSON sends status code, headers and data back to the user.
func (ar *Router) ServeJSON(w http.ResponseWriter, code int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ar.Error(w, err, http.StatusInternalServerError, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err = w.Write(data); err != nil {
		log.Printf("error writing http response: %s", err)
	}
}
