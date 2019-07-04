package adminpanel

import (
	"net/http"
)

// Setup all routes for admin panel router.
func (apr *Router) initRoutes() {
	// Do nothing on empty router (or should panic?)
	if apr.router == nil {
		return
	}

	srcHandler := http.StripPrefix("/src/", http.FileServer(http.Dir(apr.buildPath+"/src")))
	apr.router.PathPrefix("/src/").Handler(srcHandler).Methods("GET")

	managementHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, apr.buildPath+"/index.html")
	}
	apr.router.PathPrefix(`/{management:management/?}`).HandlerFunc(managementHandler).Methods("GET")

	buildHandler := http.FileServer(http.Dir(apr.buildPath))
	apr.router.PathPrefix("/").Handler(buildHandler).Methods("GET")
}
