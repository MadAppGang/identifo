package admin

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	// do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	ar.router.Path(`/{me:me/?}`).Handler(negroni.New(
		negroni.WrapFunc(ar.IsLoggedIn()),
	)).Methods("GET")

	ar.router.Path(`/{login:login/?}`).Handler(negroni.New(
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")

	ar.router.Path(`/{logout:logout/?}`).Handler(negroni.New(
		negroni.WrapFunc(ar.Logout()),
	)).Methods("POST")

	ar.router.Path(`/{apps:apps/?}`).Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchApps()),
	)).Methods("GET")
	ar.router.Path(`/{apps:apps/?}`).Handler(negroni.New(
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

	ar.router.Path(`/{users:users/?}`).Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchUsers()),
	)).Methods("GET")
	ar.router.Path(`/{users:users/?}`).Handler(negroni.New(
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

	ar.router.Path(`/{settings:settings/?}`).Handler(negroni.New(
		ar.Session(),
		negroni.WrapFunc(ar.FetchServerSettings()),
	)).Methods("GET")

	settings := mux.NewRouter().PathPrefix("/settings").Subrouter()
	ar.router.PathPrefix("/settings").Handler(negroni.New(
		ar.Session(),
		negroni.Wrap(settings),
	))
	settings.Path("/account").HandlerFunc(ar.FetchAccountSettings()).Methods("GET")
	settings.Path("/account").HandlerFunc(ar.AlterAccountSettings()).Methods("PATCH")
	settings.Path("/general").HandlerFunc(ar.AlterGeneralSettings()).Methods("PUT")
	settings.Path("/storage").HandlerFunc(ar.AlterStorageSettings()).Methods("PUT")
	settings.Path("/static").HandlerFunc(ar.AlterStaticFilesSettings()).Methods("PUT")
	settings.Path("/login").HandlerFunc(ar.AlterLoginSettings()).Methods("PUT")
	settings.Path("/services").HandlerFunc(ar.AlterExternalServicesSettings()).Methods("PUT")
	settings.Path("/storage/test").HandlerFunc(ar.TestDatabaseConnection()).Methods("POST")
}
