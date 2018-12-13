package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

//Router - router that handles all API request
type Router struct {
	middleware        *negroni.Negroni
	logger            *log.Logger
	router            *mux.Router
	appStorage        model.AppStorage
	userStorage       model.UserStorage
	tokenStorage      model.TokenStorage
	tokenService      model.TokenService
	emailService      model.EmailService
	oidcConfiguration *OIDCConfiguration
	jwk               *jwk
	Host              string
	WebRouterPrefix   string
}

//ServeHTTP identifo.Router protocol implementation
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		HostOption("http://localhost:8080"),
		WebRouterPrefixOption("/web"),
	}
}

//HostOption sets host value
func HostOption(host string) func(*Router) error {
	return func(r *Router) error {
		r.Host = host
		return nil
	}
}

//WebRouterPrefixOption sets web prefix host value
func WebRouterPrefixOption(prefix string) func(*Router) error {
	return func(r *Router) error {
		r.WebRouterPrefix = prefix
		return nil
	}
}

//NewRouter created and initiates new router
func NewRouter(logger *log.Logger, appStorage model.AppStorage, userStorage model.UserStorage, tokenStorage model.TokenStorage, tokenService model.TokenService, emailService model.EmailService, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		middleware:   negroni.Classic(),
		router:       mux.NewRouter(),
		appStorage:   appStorage,
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		tokenService: tokenService,
		emailService: emailService,
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	//setup default router to stdout
	if logger == nil {
		ar.logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
}

//ServeJSON send status code, headers and data and send it back to the user
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

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, err error, code int, userInfo string) {
	// errorResponse is a generic response for sending a error.
	type errorResponse struct {
		Error string `json:"error,omitempty"`
		Info  string `json:"info,omitempty"`
		Code  int    `json:"code,omitempty"`
	}

	// Log error.
	ar.logger.Printf("api error: %v (code=%d)", err, code)

	if err == nil {
		err = identifo.ErrorInternal
	}
	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = identifo.ErrorInternal
	}

	// Write generic error response.
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
