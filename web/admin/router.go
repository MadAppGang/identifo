package admin

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router is a router that handles admin requests.
type Router struct {
	server       model.Server
	middleware   *negroni.Negroni
	logger       *slog.Logger
	router       *mux.Router
	RedirectURL  string
	PathPrefix   string
	Host         *url.URL
	forceRestart chan<- bool
	originUpdate func() error
}

type RouterSettings struct {
	Server         model.Server
	LoggerSettings model.LoggerSettings
	Host           *url.URL
	Prefix         string
	Cors           *cors.Cors
	Restart        chan<- bool
	OriginUpdate   func() error
}

// NewRouter creates and initializes new admin router.
func NewRouter(settings RouterSettings) (model.Router, error) {
	ar := Router{
		server:       settings.Server,
		router:       mux.NewRouter(),
		Host:         settings.Host,
		PathPrefix:   settings.Prefix,
		forceRestart: settings.Restart,
		RedirectURL:  "/login",
		originUpdate: settings.OriginUpdate,
	}

	ar.logger = logging.NewLogger(
		settings.LoggerSettings.Format,
		settings.LoggerSettings.Admin.Level).
		With(logging.FieldComponent, logging.ComponentAdmin)

	ar.middleware = buildMiddleware(
		settings.LoggerSettings.DumpRequest,
		settings.LoggerSettings.Format,
		settings.LoggerSettings.Admin,
		settings.LoggerSettings.LogSensitiveData,
		settings.Cors)

	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
}

func buildMiddleware(
	dumpRequest bool,
	format string,
	logParams model.LoggerParams,
	logSensitiveData bool,
	corsHandler *cors.Cors,
) *negroni.Negroni {
	var handlers []negroni.Handler

	lm := middleware.NegroniHTTPLogger(
		logging.ComponentAdmin,
		format,
		logParams,
		model.HTTPLogDetailing(dumpRequest, logParams.HTTPDetailing),
		!logSensitiveData,
		"/login",
	)
	handlers = append(handlers, lm)

	if corsHandler == nil {
		corsHandler = cors.New(model.DefaultCors)
	}

	handlers = append(handlers,
		negroni.NewRecovery(),
		middleware.RemoveTrailingSlash(),
		corsHandler)

	return negroni.New(handlers...)
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
	ar.logger.Error("admin error",
		logging.FieldError, err,
		"errorCode", code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encodeErr := json.NewEncoder(w).Encode(&errorResponse{
		Error: err.Error(),
		Info:  userInfo,
		Code:  code,
	})
	if encodeErr != nil {
		ar.logger.Error("error writing http response",
			logging.FieldError, encodeErr)
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
		ar.logger.Error("error writing http response",
			logging.FieldError, err)
	}
}
