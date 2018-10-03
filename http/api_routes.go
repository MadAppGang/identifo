package http

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	headersOk = handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})
	originsOk = handlers.AllowedOrigins([]string{"http://localhost:8080"})
	methodsOk = handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS", "PUT", "DELETE"})
)

//setup all routes
func (ar *apiRouter) initRoutes() {
	//do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	//all API routes should have appID in it
	apiMiddlewares := ar.router.With(ar.DumpRequest(), ar.AppID())

	//setup root routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ping", ar.HandlePing()).Methods("GET")
	r.Path("/password/forgot").HandlerFunc(ar.ForgotPassword()).Methods("POST", "GET")

	// static files serve
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("../../")))

	//setup auth routes
	auth := mux.NewRouter().PathPrefix("/auth").Subrouter()
	r.PathPrefix("/auth").Handler(apiMiddlewares.With(
		ar.SignatureHandler(),
		negroni.Wrap(auth),
	))
	auth.Path("/login").HandlerFunc(ar.LoginWithPassword()).Methods("POST")
	auth.Path("/register").HandlerFunc(ar.RegisterWithPassword()).Methods("POST")

	auth.Path("/token").Handler(negroni.New(
		ar.Token("refresh"),
		negroni.Wrap(ar.RefreshToken()),
	)).Methods("GET")

	handler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	ar.router.UseHandler(handler)
}
