package http

import "github.com/rs/cors"

// Cors sets cors headers for the router
func Cors(options cors.Options) func(*apiRouter) error {
	return func(ar *apiRouter) error {
		return ar.setCORS(options)
	}
}

func (ar *apiRouter) setCORS(options cors.Options) error {
	c := cors.New(options)
	ar.router.Use(c)
	return nil
}
