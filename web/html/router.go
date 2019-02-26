package html

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

//Router routes all html files and handles html post requests
type Router struct {
	Middleware      *negroni.Negroni
	Logger          *log.Logger
	Router          *mux.Router
	AppStorage      model.AppStorage
	UserStorage     model.UserStorage
	TokenStorage    model.TokenStorage
	TokenService    model.TokenService
	EmailService    model.EmailService
	StaticPages     StaticPages
	StaticFilesPath StaticFilesPath
	EmailTemplates  EmailTemplates
	PathPrefix      string
	Host            string
}

func defaultOptions() []func(*Router) error {
	return []func(*Router) error{
		DefaultStaticPagesOptions(),
		DefaultStaticPathOptions(),
		PathPrefixOptions("/web"),
		HostOption("http://localhost:8080"),
	}
}

// PathPrefixOptions set path prefix options
func PathPrefixOptions(prefix string) func(r *Router) error {
	return func(r *Router) error {
		r.PathPrefix = prefix
		return nil
	}
}

// HostOption set hostname
func HostOption(host string) func(r *Router) error {
	return func(r *Router) error {
		r.Host = host
		return nil
	}
}

//NewRouter created and initiates new router
func NewRouter(logger *log.Logger, appStorage model.AppStorage, userStorage model.UserStorage, tokenStorage model.TokenStorage, tokenService model.TokenService, emailService model.EmailService, options ...func(*Router) error) (model.Router, error) {
	ar := Router{
		Middleware:   negroni.Classic(),
		Router:       mux.NewRouter(),
		AppStorage:   appStorage,
		UserStorage:  userStorage,
		TokenStorage: tokenStorage,
		TokenService: tokenService,
		EmailService: emailService,
	}

	for _, option := range append(defaultOptions(), options...) {
		if err := option(&ar); err != nil {
			return nil, err
		}
	}

	//setup default router to stdout
	if logger == nil {
		ar.Logger = log.New(os.Stdout, "HTML_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.initRoutes()
	ar.Middleware.UseHandler(ar.Router)
	return &ar, nil
}

// Error writes an API error message to the response and logger.
func (ar *Router) Error(w http.ResponseWriter, err error, code int, userInfo string) {
	// Log error.
	ar.Logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = identifo.ErrorInternal
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseString := `
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Home Network</title>
	</head>
	<body>
	<h2>Error</h2></br>
	<h3>
	` +
		fmt.Sprintf("Error: %s, code: %d, userInfo: %s", err.Error(), code, userInfo) +
		`
	</h3>
	</body>
	</html>
	`
	w.WriteHeader(code)
	if _, wrErr := io.WriteString(w, responseString); wrErr != nil {
		ar.Logger.Println("Error writing response string:", wrErr)
	}
}

//ServeHTTP identifo.Router protocol implementation
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.Router.ServeHTTP(w, r)
}
