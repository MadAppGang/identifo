package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"runtime"
	"sync"

	"github.com/gorilla/mux"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/rs/cors"
)

// Router is a router that handles all API requests.
type Router struct {
	server             model.Server
	cors               *cors.Cors
	logger             *slog.Logger
	router             *mux.Router
	tfaType            model.TFAType
	tfaResendTimeout   int
	oidcConfiguration  *OIDCConfiguration
	jwk                *JWK
	Authorizer         *authorization.Authorizer
	Host               *url.URL
	SupportedLoginWays model.LoginWith

	tokenPayloadServices     map[string]model.TokenPayloadProvider
	tokenPayloadServicesLock sync.RWMutex

	ls *l.Printer // localized string
}

type RouterSettings struct {
	Server           model.Server
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
		server:               settings.Server,
		router:               mux.NewRouter(),
		Authorizer:           settings.Authorizer,
		Host:                 settings.Host,
		tfaType:              settings.TFAType,
		tfaResendTimeout:     settings.TFAResendTimeout,
		SupportedLoginWays:   settings.LoginWith,
		cors:                 settings.Cors,
		ls:                   l,
		tokenPayloadServices: map[string]model.TokenPayloadProvider{},
	}

	ar.logger = logging.NewLogger(
		settings.LoggerSettings.Format,
		settings.LoggerSettings.API.Level).
		With(logging.FieldComponent, logging.ComponentAPI)

	ar.initRoutes(settings.LoggerSettings)

	return &ar, nil
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}

// ServeJSON sends status code, headers and data and send it back to the user
func (ar *Router) ServeJSON(w http.ResponseWriter, locale string, status int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(data); err != nil {
		ar.logger.Error("error writing http response",
			logging.FieldError, err)
	}
}

func NewLocalizedError(status int, locale string, errID l.LocalizedString, details ...any) *LocalizedError {
	return &LocalizedError{
		status:  status,
		locale:  locale,
		errID:   errID,
		details: details,
	}
}

type LocalizedError struct {
	status  int
	locale  string
	errID   l.LocalizedString
	details []any
}

func (e *LocalizedError) Error() string {
	return fmt.Sprintf("localized error: %v (status=%v). Details: %v.", e.errID, e.status, e.details)
}

func (ar *Router) ErrorResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *LocalizedError:
		ar.error(w, 3, e.locale, e.status, e.errID, e.details...)
	default:
		ar.error(w, 3, "", http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
	}
}

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, locale string, status int, errID l.LocalizedString, details ...any) {
	ar.error(w, 2, locale, status, errID, details...)
}

func (ar *Router) error(w http.ResponseWriter, callerDepth int, locale string, status int, errID l.LocalizedString, details ...any) {
	// errorResponse is a generic response for sending an error.
	type errorResponse struct {
		ID       string `json:"id"`
		Message  string `json:"message,omitempty"`
		Status   int    `json:"status"`
		Location string `json:"location"`
	}

	if errID == "" {
		errID = l.APIInternalServerError
	}

	_, file, no, ok := runtime.Caller(callerDepth)
	if !ok {
		file = "unknown file"
		no = -1
	}
	message := ar.ls.SL(locale, errID, details...)

	logLevel := slog.LevelWarn
	if status >= 500 {
		logLevel = slog.LevelError
	}

	ar.logger.Log(
		context.Background(),
		logLevel,
		"api error",
		logging.FieldErrorID, errID,
		"status", status,
		"details", message,
		"where", fmt.Sprintf("%v:%d", file, no),
	)

	// Write generic error response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := struct {
		Error errorResponse `json:"error"`
	}{
		Error: errorResponse{
			ID:       string(errID),
			Message:  message,
			Status:   status,
			Location: fmt.Sprintf("%s:%d", file, no),
		},
	}

	encodeErr := json.NewEncoder(w).Encode(resp)
	if encodeErr != nil {
		ar.logger.Error("error writing http response",
			logging.FieldError, encodeErr)
	}
}
