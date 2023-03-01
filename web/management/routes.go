package management

import (
	"time"

	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func logger(name string) {
}

// setup all routes
func (ar *Router) initRoutes() {
	// A good base middleware stack
	logger := httplog.LoggerWithName("MANAGEMENT")

	if ar.loggerSettings.DumpRequest {
		logger = httplog.LoggerWithFormatterAndName("MANAGEMENT", httplog.DefaultLogFormatterWithRequestHeadersAndBody)
	}

	ar.router.Use(middleware.RequestID)
	ar.router.Use(middleware.RealIP)
	ar.router.Use(logger)
	ar.router.Use(middleware.Recoverer)
	ar.router.Use(middleware.CleanPath)
	ar.router.Use(middleware.Timeout(30 * time.Second))
	ar.router.Use(ar.AuthMiddleware)

	ar.router.Get("/test", ar.test)
	// token endpoints
	ar.router.Route("/token", func(r chi.Router) {
		r.Post("/invite", ar.getInviteToken)
		r.Post("/reset_password", ar.getResetPasswordToken)
	})
}
