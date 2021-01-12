package model

import (
	"net/http"

	"github.com/madappgang/identifo/plugin/shared"
)

// Server holds together all dependencies.
type Server interface {
	Router() Router
	AppStorage() AppStorage
	UserStorage() shared.UserStorage
	ConfigurationStorage() ConfigurationStorage
	ImportApps(filename string) error
	ImportUsers(filename string) error
	Close()
}

// Router handles all incoming http requests.
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
