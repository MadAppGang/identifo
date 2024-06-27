package management

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	imiddleware "github.com/madappgang/identifo/v2/web/middleware"
)

// setup all routes
func (ar *Router) initRoutes(loggerSettings model.LoggerSettings) {
	// A good base middleware stack
	lm := imiddleware.HTTPLogger(
		logging.ComponentAPI,
		loggerSettings.Format,
		loggerSettings.Management,
		model.HTTPLogDetailing(loggerSettings.DumpRequest, loggerSettings.Management.HTTPDetailing),
	)

	ar.router.Use(middleware.RequestID)
	ar.router.Use(middleware.RealIP)
	ar.router.Use(lm)
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
