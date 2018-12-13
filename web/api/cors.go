package api

import "github.com/rs/cors"

// Cors sets cors headers for the router
func Cors(options cors.Options) func(*Router) error {
	return func(ar *Router) error {
		return ar.setCORS(options)
	}
}

func (ar *Router) setCORS(options cors.Options) error {
	c := cors.New(options)
	ar.middleware.Use(c)
	return nil
}
