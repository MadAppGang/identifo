package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/madappgang/identifo"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

//apiRoutes - router that handles all API request
type apiRouter struct {
	router       *negroni.Negroni
	logger       *log.Logger
	appStorage   model.AppStorage
	userStorage  model.UserStorage
	tokenStorage model.TokenStorage
	tokenService model.TokenService
}

//ServeHTTP identifo.Router protocol implementation
func (ar *apiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}

//NewRouter created and initiates new router
func NewRouter(logger *log.Logger, appStorage model.AppStorage, userStorage model.UserStorage, tokenStorage model.TokenStorage) model.Router {
	ar := apiRouter{}
	ar.router = negroni.Classic()
	//setup default router to stdout
	if logger == nil {
		ar.logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	ar.appStorage = appStorage
	ar.userStorage = userStorage
	ar.tokenStorage = tokenStorage
	return &ar
}

//ServeJSON send status code, headers and data and send it back to the user
func (ar *apiRouter) ServeJSON(w http.ResponseWriter, code int, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		ar.Error(w, err, http.StatusInternalServerError, "")
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&errorResponse{
		Error: err.Error(),
		Info:  userInfo,
		Code:  code,
	})
}
