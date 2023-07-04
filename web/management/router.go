package management

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

type RouterSettings struct {
	Server             model.Server
	Logger             *log.Logger
	LoggerSettings     model.LoggerSettings
	Storage            model.ManagementKeysStorage
	Locale             string
	SupportedLoginWays model.LoginWith
}

type Router struct {
	server         model.Server
	ls             *l.Printer // localized string
	logger         *log.Logger
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

	// setup logger to stdout.
	if settings.Logger == nil {
		ar.logger = log.New(os.Stdout, "MANAGEMENT_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		ar.logger = settings.Logger
	}

	ar.initRoutes()

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
		log.Printf("error writing http response: %s", err)
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

	// Log error.
	ar.logger.Printf("api error: %v (status=%v). Details: %v. Where: %v:%d.", errID, status, message, file, no)

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
		ar.logger.Printf("error writing http response: %s", errID)
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
