package model

import (
	"net/http"
)

//Server holds together all dependencies
type Server struct {
	Router Router
}

//Router is class to handle all incoming http requests
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
