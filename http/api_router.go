package http

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

//apiRoutes - router that handles all API request
type apiRouter struct {
	router            *negroni.Negroni
	logger            *log.Logger
	appStorage        model.AppStorage
	userStorage       model.UserStorage
	tokenStorage      model.TokenStorage
	tokenService      model.TokenService
	handler           *mux.Router
	oidcConfiguration *OIDCConfiguration
	jwk               *jwk
}

//ServeHTTP identifo.Router protocol implementation
func (ar *apiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}

func defaultOptions() []func(*apiRouter) error {
	return []func(*apiRouter) error{ServeDefaultStaticPages()}
}

//NewRouter created and initiates new router
func NewRouter(logger *log.Logger, appStorage model.AppStorage, userStorage model.UserStorage, tokenStorage model.TokenStorage, tokenService model.TokenService, options ...func(*apiRouter) error) (model.Router, error) {
	ar := apiRouter{
		router:       negroni.Classic(),
		handler:      mux.NewRouter(),
		appStorage:   appStorage,
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		tokenService: tokenService,
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
	ar.router.UseHandler(ar.handler)

	return &ar, nil
}

//ServeJSON send status code, headers and data and send it back to the user
func (ar *apiRouter) ServeJSON(w http.ResponseWriter, code int, v interface{}) {
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
func (ar *apiRouter) Error(w http.ResponseWriter, err error, code int, userInfo string) {
	// errorResponse is a generic response for sending a error.
	type errorResponse struct {
		Error string `json:"error,omitempty"`
		Info  string `json:"info,omitempty"`
		Code  int    `json:"code,omitempty"`
	}

	// Log error.
	ar.logger.Printf("http error: %s (code=%d)", err, code)

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
