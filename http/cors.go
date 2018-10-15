package http

import "github.com/rs/cors"

func (ar *apiRouter) initCORS(o cors.Options) {
	c := cors.New(o)
	ar.router.Use(c)
}
