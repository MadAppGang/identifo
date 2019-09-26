package model

import (
	"github.com/rs/cors"
)

// CorsOptions are options for routers CORS.
type CorsOptions struct {
	Admin *cors.Options
	API   *cors.Options
	HTML  *cors.Options
}
