package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	if ar.router == nil {
		panic("Empty admin router")
	}

	ar.router.Path("/me").Handler(negroni.New(
		negroni.WrapFunc(ar.IsLoggedIn()),
	)).Methods("GET")

	ar.router.Path("/login").Handler(negroni.New(
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")

	ar.router.Path("/logout").Handler(negroni.New(
		negroni.WrapFunc(ar.Logout()),
	)).Methods("POST")

	ar.router.Path("/apps").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchApps()),
	)).Methods("GET")
	ar.router.Path("/federated-providers").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FederatedProvidersList()),
	)).Methods("GET")
	ar.router.Path("/apps").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.CreateApp()),
	)).Methods("POST")

	apps := mux.NewRouter().PathPrefix("/apps").Subrouter()
	ar.router.PathPrefix("/apps").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(apps),
	))
	apps.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.GetApp()).Methods("GET")
	apps.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.UpdateApp()).Methods("PUT")
	apps.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.DeleteApp()).Methods("DELETE")
	apps.Path("/").HandlerFunc(ar.DeleteAllApps()).Methods("DELETE")

	ar.router.Path("/users").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchUsers()),
	)).Methods("GET")
	ar.router.Path("/users").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.CreateUser()),
	)).Methods("POST")

	users := mux.NewRouter().PathPrefix("/users").Subrouter()
	ar.router.PathPrefix("/users").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(users),
	))
	users.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.GetUser()).Methods("GET")
	users.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.UpdateUser()).Methods("PUT")
	users.Path("/{id:[a-zA-Z0-9]+}").HandlerFunc(ar.DeleteUser()).Methods("DELETE")
	users.Path("/generate_new_reset_token").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.GenerateNewResetTokenUser()),
	)).Methods("POST")

	ar.router.Path("/settings").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchSettings()),
	)).Methods("GET")

	ar.router.Path("/settings").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.UpdateSettings()),
	)).Methods("PUT")

	ar.router.Path("/test_connection").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.TestConnection()),
	)).Methods("POST")

	ar.router.Path("/generate_new_secret").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.GenerateNewSecret()),
	)).Methods("POST")

	ar.router.Path("/invites").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchInvites()),
	)).Methods("GET")

	ar.router.Path("/invites").Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.AddInvite()),
	)).Methods("POST")

	invites := mux.NewRouter().PathPrefix("/invites").Subrouter()
	ar.router.PathPrefix("/invites").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(invites),
	))

	invites.Path("{id:[a-zA-Z0-9]+}").HandlerFunc(ar.GetInviteByID()).Methods(http.MethodGet)
	invites.Path("{id:[a-zA-Z0-9]+}").HandlerFunc(ar.ArchiveInviteByID()).Methods(http.MethodDelete)

	static := mux.NewRouter().PathPrefix("/static").Subrouter()
	ar.router.PathPrefix("/static").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(static),
	))

	static.Path("/uploads/keys").HandlerFunc(ar.UploadJWTKeys()).Methods("POST")
	static.Path("/keys").HandlerFunc(ar.GetJWTKeys()).Methods("GET")
}
