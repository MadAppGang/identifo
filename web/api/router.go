package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

// Router is a router that handles all API requests.
type Router struct {
	middleware              *negroni.Negroni
	logger                  *log.Logger
	router                  *mux.Router
	appStorage              model.AppStorage
	userStorage             model.UserStorage
	tokenStorage            model.TokenStorage
	verificationCodeStorage model.VerificationCodeStorage
	tokenService            jwtService.TokenService
	smsService              model.SMSService
	emailService            model.EmailService
	oidcConfiguration       *OIDCConfiguration
	jwk                     *jwk
	Host                    string
	WebRouterPrefix         string
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.router.ServeHTTP(w, r)
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		WebRouterPrefixOption("/web"),
	}
}

// HostOption sets host value.
func HostOption(host string) func(*Router) error {
	return func(r *Router) error {
		r.Host = host
		return nil
	}
}

// WebRouterPrefixOption sets web prefix host value.
func WebRouterPrefixOption(prefix string) func(*Router) error {
	return func(r *Router) error {
		r.WebRouterPrefix = prefix
		return nil
	}
}

// NewRouter creates and initilizes new router.
func NewRouter(logger *log.Logger, appStorage model.AppStorage, userStorage model.UserStorage, tokenStorage model.TokenStorage, verificationCodeStorage model.VerificationCodeStorage, tokenService jwtService.TokenService, smsService model.SMSService, emailService model.EmailService, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		middleware:              negroni.Classic(),
		router:                  mux.NewRouter(),
		appStorage:              appStorage,
		userStorage:             userStorage,
		tokenStorage:            tokenStorage,
		verificationCodeStorage: verificationCodeStorage,
		tokenService:            tokenService,
		smsService:              smsService,
		emailService:            emailService,
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	// setup logger to stdout.
	if logger == nil {
		ar.logger = log.New(os.Stdout, "API_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.initRoutes()
	ar.middleware.UseHandler(ar.router)

	return &ar, nil
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
	// errorResponse is a generic response for sending a error.
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
		Message:         GetMessage(errID),
		DetailedMessage: details,
		Status:          status,
	}})
	if encodeErr != nil {
		ar.logger.Printf("error writing http response: %s", errID)
	}
}
