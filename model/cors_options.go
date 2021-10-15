package model

import (
	"github.com/rs/cors"
)

var DefaultCors = cors.Options{AllowedHeaders: []string{"*", "x-identifo-clientid"}, AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}, AllowCredentials: true}
