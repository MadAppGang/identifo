package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Router is a router that handles all API requests.
type Router struct {
	server               model.Server
	middleware           *negroni.Negroni
	cors                 *cors.Cors
	logger               *log.Logger
	router               *mux.Router
	tfaType              model.TFAType
	tfaResendTimeout     int
	oidcConfiguration    *OIDCConfiguration
	jwk                  *jwk
	Authorizer           *authorization.Authorizer
	Host                 *url.URL
	SupportedLoginWays   model.LoginWith
	tokenPayloadServices map[string]model.TokenPayloadProvider
	LoggerSettings       model.LoggerSettings
	ls                   *message.Printer // localized string
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
	l, err := language.Parse(settings.Locale)
	// go to default english locale
	if err != nil {
		l = language.English
	}
	p := message.NewPrinter(l)

	ar := Router{
		server:             settings.Server,
		middleware:         negroni.New(middleware.NewNegroniLogger("API"), negroni.NewRecovery()),
		router:             mux.NewRouter(),
		Authorizer:         settings.Authorizer,
		LoggerSettings:     settings.LoggerSettings,
		Host:               settings.Host,
		tfaType:            settings.TFAType,
		tfaResendTimeout:   settings.TFAResendTimeout,
		SupportedLoginWays: settings.LoginWith,
		cors:               settings.Cors,
		ls:                 p,
	}

	// setup logger to stdout.
	if settings.Logger == nil {
		ar.logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		ar.logger = settings.Logger
	}

	ar.tokenPayloadServices = make(map[string]model.TokenPayloadProvider)

	ar.middleware.Use(ar.RemoveTrailingSlash())

	if ar.cors != nil {
		ar.middleware.Use(ar.cors)
	}
	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reroute to our internal implementation
	ar.middleware.ServeHTTP(w, r)
}

// ServeJSON sends status code, headers and data and send it back to the user
func (ar *Router) ServeJSON(w http.ResponseWriter, status int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "Unable to marshall response. Err: "+err.Error(), "Router.ServerJSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(data); err != nil {
		log.Printf("error writing http response: %s", err)
	}
}

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, errID MessageID, status int, details, where string) {
	// errorResponse is a generic response for sending an error.
	type errorResponse struct {
		ID              MessageID `json:"id"`
		Message         string    `json:"message,omitempty"`
		DetailedMessage string    `json:"detailed_message,omitempty"`
		Status          int       `json:"status"`
	}

	// Log error.
	ar.logger.Printf("api error: %v (status=%v). Details: %v. Where: %v.", errID, status, details, where)

	if errID == "" {
		errID = ErrorAPIInternalServerError
	}
	// Hide error from client if it is internal.
	if status == http.StatusInternalServerError {
		errID = ErrorAPIInternalServerError
	}

	// Write generic error response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encodeErr := json.NewEncoder(w).Encode(map[string]interface{}{"error": &errorResponse{
		ID:              errID,
		Message:         ar.ls.Sprintf(GetMessage(errID)),
		DetailedMessage: ar.ls.Sprintf(details),
		Status:          status,
	}})
	if encodeErr != nil {
		ar.logger.Printf("error writing http response: %s", errID)
	}
}
