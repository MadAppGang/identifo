package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

// Router is a router that handles admin requests.
type Router struct {
	middleware     *negroni.Negroni
	logger         *log.Logger
	router         *mux.Router
	sessionService model.SessionService
	sessionStorage model.SessionStorage
	userStorage    model.UserStorage
	ConfigPath     string
	RedirectURL    string
	PathPrefix     string
	Host           string
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		PathPrefixOptions("/admin"),
		HostOption("http://localhost:8080"),
		RedirectURLOption("/login"),
	}
}

// HostOption sets host value.
func HostOption(host string) func(*Router) error {
	return func(r *Router) error {
		r.Host = host
		return nil
	}
}

// ConfigPathOption sets path to configuration file.
func ConfigPathOption(configPath string) func(*Router) error {
	return func(r *Router) error {
		r.ConfigPath = configPath
		return nil
	}
}

// RedirectURLOption sets redirect url value.
func RedirectURLOption(redirectURL string) func(*Router) error {
	return func(r *Router) error {
		r.RedirectURL = path.Join(r.Host, r.PathPrefix, redirectURL)
		return nil
	}
}

// PathPrefixOptions sets path prefix options.
func PathPrefixOptions(prefix string) func(r *Router) error {
	return func(r *Router) error {
		r.PathPrefix = prefix
		return nil
	}
}

// NewRouter creates and initializes new admin router.
func NewRouter(logger *log.Logger, sessionService model.SessionService, sessionStorage model.SessionStorage, userStorage model.UserStorage, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		middleware:     negroni.Classic(),
		router:         mux.NewRouter(),
		sessionService: sessionService,
		sessionStorage: sessionStorage,
		userStorage:    userStorage,
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	if logger == nil {
		ar.logger = log.New(os.Stdout, "ADMIN_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

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

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = identifo.ErrorInternal
	}

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
	//reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
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
