package adminpanel

import (
	"log"
	"net/http"
	"os"
)

// Setup all routes for admin panel router.
func (apr *Router) initRoutes() {
	// Do nothing on empty router (or should panic?)
	if apr.router == nil {
		return
	}
	f, err := os.Getwd()
	log.Println(apr.buildPath, f, err)

	buildHandler := http.FileServer(http.Dir("./admin_panel/build"))
	apr.router.PathPrefix("/").Handler(buildHandler).Methods("GET")
}
