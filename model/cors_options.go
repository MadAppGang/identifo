package model

import (
	"github.com/go-chi/cors"
)

var DefaultCors = cors.Options{
	AllowedHeaders:   []string{"*", "x-identifo-clientid"},
	AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
	AllowCredentials: true,
}
