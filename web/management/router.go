package management

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

type RouterSettings struct {
	Server             model.Server
	LoggerSettings     model.LoggerSettings
	Storage            model.ManagementKeysStorage
	Locale             string
	SupportedLoginWays model.LoginWith
}

type Router struct {
	server         model.Server
	ls             *l.Printer // localized string
	logger         *slog.Logger
	router         *chi.Mux
	loggerSettings model.LoggerSettings
	stor           model.ManagementKeysStorage
	loginWith      model.LoginWith
}

// NewRouter creates and inits new router.
func NewRouter(settings RouterSettings) (*Router, error) {
	l, err := l.NewPrinter(settings.Locale)
	if err != nil {
		return nil, err
	}

	ar := Router{
		server:         settings.Server,
		router:         chi.NewRouter(),
		ls:             l,
		loggerSettings: settings.LoggerSettings,
		stor:           settings.Storage,
		loginWith:      settings.SupportedLoginWays,
	}

	ar.logger = logging.NewLogger(
		settings.LoggerSettings.Format,
		settings.LoggerSettings.Management.Level,
	).With(logging.FieldComponent, logging.ComponentManagement)

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
		logging.DefaultLogger.Error("error writing http response", logging.FieldError, err)
	}
}

var jsonOkBody = []byte(`{"result": "ok"}`)

func (ar *Router) ServeJSONOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonOkBody); err != nil {
		ar.logger.Error("error writing http response",
			logging.FieldError, err)
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

	ar.logger.Error("api error",
		logging.FieldErrorID, errID,
		"status", status,
		"details", message,
		"where", fmt.Sprintf("%v:%d", file, no))

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

// MustParseJSON parses request body json data to the `out` struct.
// If error happens, writes it to ResponseWriter.
func (ar *Router) MustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	locale := r.Header.Get("Accept-Language")

	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		return err
	}

	return nil
}
