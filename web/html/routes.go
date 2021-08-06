package html

import (
	"fmt"
	"net/http"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	if ar.Router == nil {
		panic("Empty HTML router")
	}

	appHandler := ar.Server.Storages().Static.WebHandlers()
	ar.Router.PathPrefix(`/`).Handler(appHandler.AppHandler).Methods("GET")
	ar.Router.PathPrefix(`/`).HandlerFunc(ar.AppPost()).Methods("POST")
}

func (ar *Router) AppPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if len(r.Form) > 0 {
			r.URL.RawQuery = r.Form.Encode()
			fmt.Println(r.URL.String())
			fmt.Println(r.URL.Path)
			http.Redirect(w, r, ar.PathPrefix+r.URL.String(), http.StatusMovedPermanently)
		} else {
			appHandler := ar.Server.Storages().Static.WebHandlers()
			appHandler.AppHandler.ServeHTTP(w, r)
		}
	}
}
