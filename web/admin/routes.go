package admin

import (
	"github.com/urfave/negroni"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	// do nothing on empty router (or should panic?)
	if ar.router == nil {
		return
	}

	ar.router.Path("/{login:login\\/?}").Handler(negroni.New(
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")
}
