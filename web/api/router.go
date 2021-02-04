package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/utils/originchecker"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Router is a router that handles all API requests.
type Router struct {
	middleware              *negroni.Negroni
	cors                    *cors.Cors
	logger                  *log.Logger
	router                  *mux.Router
	appStorage              model.AppStorage
	userStorage             model.UserStorage
	tokenStorage            model.TokenStorage
	tokenBlacklist          model.TokenBlacklist
	verificationCodeStorage model.VerificationCodeStorage
	staticFilesStorage      model.StaticFilesStorage
	tfaType                 model.TFAType
	tokenService            jwtService.TokenService
	smsService              model.SMSService
	emailService            model.EmailService
	oidcConfiguration       *OIDCConfiguration
	jwk                     *jwk
	Authorizer              *authorization.Authorizer
	Host                    string
	SupportedLoginWays      model.LoginWith
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

// CorsOption sets cors option.
func CorsOption(corsOptions *model.CorsOptions, originChecker *originchecker.OriginChecker) func(*Router) error {
	return func(r *Router) error {
		if corsOptions != nil && corsOptions.API != nil {
			if originChecker != nil {
				corsOptions.API.AllowOriginRequestFunc = originChecker.With(corsOptions.API.AllowOriginRequestFunc).CheckOrigin
			}
			r.cors = cors.New(*corsOptions.API)
		}
		return nil
	}
}

// SupportedLoginWaysOption is for setting supported ways of logging in into the app.
func SupportedLoginWaysOption(loginWays model.LoginWith) func(*Router) error {
	return func(r *Router) error {
		r.SupportedLoginWays = loginWays
		return nil
	}
}

// TFATypeOption is for setting two-factor authentication type.
func TFATypeOption(tfaType model.TFAType) func(*Router) error {
	return func(r *Router) error {
		r.tfaType = tfaType
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
func NewRouter(logger *log.Logger, as model.AppStorage, us model.UserStorage, ts model.TokenStorage, tb model.TokenBlacklist, vcs model.VerificationCodeStorage, sfs model.StaticFilesStorage, tServ jwtService.TokenService, smsServ model.SMSService, emailServ model.EmailService, authorizer *authorization.Authorizer, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		middleware:              negroni.Classic(),
		router:                  mux.NewRouter(),
		appStorage:              as,
		userStorage:             us,
		tokenStorage:            ts,
		tokenBlacklist:          tb,
		verificationCodeStorage: vcs,
		staticFilesStorage:      sfs,
		tokenService:            tServ,
		smsService:              smsServ,
		emailService:            emailServ,
		Authorizer:              authorizer,
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

	if ar.cors != nil {
		ar.middleware.Use(ar.cors)
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
	// errorResponse is a generic response for sending an error.
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
