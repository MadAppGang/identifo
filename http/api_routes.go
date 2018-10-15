package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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

	//static pages
	r.HandleFunc("/login", ar.ServeTemplate(ar.staticPages.Login)).Methods("GET")
	r.HandleFunc("/registration", ar.ServeTemplate(ar.staticPages.Registration)).Methods("GET")
	r.HandleFunc("/password/forgot", ar.ServeTemplate(ar.staticPages.ForgotPassword)).Methods("GET")
	r.HandleFunc("/password/reset", ar.ServeTemplate(ar.staticPages.ResetPassword)).Methods("GET")

	//static files
	handler := http.FileServer(http.Dir("../../static"))
	r.PathPrefix("/css").Handler(handler)
	r.PathPrefix("/js").Handler(handler)

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

	ar.router.UseHandler(r)
}
