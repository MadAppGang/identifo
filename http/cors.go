package http

import "github.com/rs/cors"

func (ar *apiRouter) AddCORS(options cors.Options) {
	c := cors.New(options)
	ar.router.Use(c)
}
