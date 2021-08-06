package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/utils/originchecker"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router is a router that handles admin requests.
type Router struct {
	server        model.Server
	middleware    *negroni.Negroni
	cors          *cors.Cors
	originChecker *originchecker.OriginChecker
	logger        *log.Logger
	router        *mux.Router
	newSettings   *model.ServerSettings
	RedirectURL   string
	PathPrefix    string
	Host          string
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		PathPrefixOptions("/admin"),
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

// CorsOption sets cors option.
func CorsOption(corsOptions model.CorsOptions, originChecker *originchecker.OriginChecker) func(*Router) error {
	return func(r *Router) error {
		if originChecker != nil {
			r.originChecker = originChecker
		} else {
			r.originChecker = originchecker.NewOriginChecker()
		}

		if corsOptions.Admin != nil {
			r.cors = cors.New(*corsOptions.Admin)
		}
		return nil
	}
}

// // ServerConfigPathOption sets path to configuration file with admin server settings.
// func ServerConfigPathOption(configPath string) func(*Router) error {
// 	return func(r *Router) error {
// 		r.ServerConfigPath = configPath
// 		return nil
// 	}
// }

// // ServerSettingsOption sets path to configuration file with server settings.
// func ServerSettingsOption(settings *model.ServerSettings) func(*Router) error {
// 	return func(r *Router) error {
// 		r.ServerSettings = settings
// 		r.newSettings = settings
// 		return nil
// 	}
// }

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
func NewRouter(server model.Server, logger *log.Logger, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		server:     server,
		middleware: negroni.Classic(),
		router:     mux.NewRouter(),
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	if logger == nil {
		ar.logger = log.New(os.Stdout, "ADMIN_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.middleware.Use(ar.RemoveTrailingSlash())

	if ar.cors == nil {
		ar.cors = ar.defaultCORS()
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

func (ar *Router) defaultCORS() *cors.Cors {
	allowedOrigins := []string{"http://localhost:*"}
	if adminPanelURL := os.Getenv("ADMIN_PANEL_URL"); len(adminPanelURL) > 0 {
		allowedOrigins = append(allowedOrigins, adminPanelURL)
	}

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	})
}
