package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/madappgang/identifo/v2/model"
	// wm "github.com/madappgang/identifo/v2/web/middleware"
)

// setup all routes
// inspired by https://auth0.com/docs/api/authentication#introduction
// https://reference.clerk.dev/reference/frontend-api-reference/users/introduction

func (ar *Router) initRoutes() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(ar.cors))
	r.Use(ar.ConfigCheck)

	r.Get("/ping", ar.Ping)
	r.Get("/hello", ar.Hello)

	// auth router for non authenticated users
	r.Route("/auth", func(r chi.Router) {
		r.Use(ar.HMACSignature)
		r.Route("/passwordless", func(r chi.Router) {
			r.Post("/start", ar.RequestChallenge())
			r.Post("/complete", ar.PasswordlessLogin())
		})
		r.Post("/login", ar.LoginWithPassword())
		r.Post("/register", ar.RegisterWithPassword())
		r.Route("/password", func(r chi.Router) {
			r.Post("/reset", ar.RequestResetPassword())
			// the only route here for authenticated user, but he have to use reset bearer token, not access
			r.With(ar.Token(model.TokenTypeReset)).Post("/change", ar.ChangePassword())
		})
		r.Route("/federated", func(r chi.Router) {
			// r.Handle("/start", ar.FederatedStart())
			// r.Handle("/complete", ar.FederatedComplete())
		})
		r.With(ar.Token(model.TokenTypeRefresh)).Post("/token", ar.RefreshTokens())
	})

	// routes for apps
	r.Route("/app", func(r chi.Router) {
		r.Use(ar.HMACSignature)
		r.Get("/settings", ar.GetAppSettings())
	})

	// authenticated user routes
	r.Route("/user", func(r chi.Router) {
		r.Use(ar.Token(model.TokenTypeAccess))
		r.Get("/profile", ar.GetUser())
		r.Get("/profile/{id}", ar.GetUser())
		r.Put("/profile", ar.UpdateUser())
		r.Patch("/profile", ar.UpdateUser())
		// TODO
		// r.Get("/data", ar.GetUserData())
		// r.Put("/data", ar.UpdateUserData())
		r.Route("/device", func(r chi.Router) {
			// TODO: add device management endpoints: register, get, list, update
		})
		r.Post("/invite", ar.RequestInviteLink())
		r.Post("/logout", ar.Logout())
	})

	// introspection routes
	r.Route("/.well-known", func(r chi.Router) {
		r.Get("/openid-configuration", ar.OIDCConfiguration())
		r.Get("/jwks.json", ar.OIDCJwks())
	})
	ar.router = r
}
