package model

import (
	"net/http"
)

// Server holds together all dependencies.
type Server interface {
	Router() Router
	AppStorage() AppStorage
	UserStorage() UserStorage
	ImportApps(filename string) error
	ImportUsers(filename string) error
}

// Router handles all incoming http requests.
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
