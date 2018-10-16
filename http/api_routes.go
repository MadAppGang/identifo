package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//setup all routes
func (ar *apiRouter) initRoutes(staticPages *StaticPages) {
	//do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	//all API routes should have appID in it
	apiMiddlewares := ar.router.With(ar.DumpRequest(), ar.AppID())

	//setup root routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ping", ar.HandlePing()).Methods("GET")

	//setup routes for static pages
	if staticPages != nil {
		static := r.NewRoute().Subrouter()

		static.HandleFunc("/login", ar.ServeTemplate(staticPages.Login)).Methods("GET")
		static.HandleFunc("/register", ar.ServeTemplate(staticPages.Registration)).Methods("GET")
		static.HandleFunc("/password/forgot", ar.ServeTemplate(staticPages.ForgotPassword)).Methods("GET")
		static.HandleFunc("/password/reset", ar.ServeTemplate(staticPages.ResetPassword)).Methods("GET")

		//setup routes for static files
		handler := http.FileServer(http.Dir("../../static"))
		static.PathPrefix("/css").Handler(handler)
		static.PathPrefix("/js").Handler(handler)

		r.NewRoute().Handler(static)
	}

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
