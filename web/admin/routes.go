package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/madappgang/identifo/v2/model"
	wm "github.com/madappgang/identifo/v2/web/middleware"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(ar.cors))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", ar.Login())
		r.With(ar.token()).Post("/logout", ar.Logout())
	})

	r.Route("/apps", func(r chi.Router) {
		r.Use(ar.token())
		r.Get("/", ar.FetchApps())
		r.Post("/", ar.CreateApp())
		r.Get("/{id:[a-zA-Z0-9]+}", ar.GetApp())
		r.Put("/{id:[a-zA-Z0-9]+}", ar.UpdateApp())
		r.Delete("/{id:[a-zA-Z0-9]+}", ar.DeleteApp())
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(ar.token())
		r.Get("/", ar.FetchUsers())
		r.Post("/", ar.CreateUser())
		r.Get("/{id:[a-zA-Z0-9]+}", ar.GetUser())
		r.Put("/{id:[a-zA-Z0-9]+}", ar.UpdateUser())
		r.Delete("/{id:[a-zA-Z0-9]+}", ar.DeleteUser())
		r.Post("/{id:[a-zA-Z0-9]+}/reset", ar.GenerateNewResetTokenUser())
	})

	r.Route("/settings", func(r chi.Router) {
		r.Use(ar.token())
		r.Get("/", ar.FetchSettings())
		r.Put("/", ar.UpdateSettings())
		r.Post("/test", ar.TestConnection())
		r.Post("/new_secret", ar.GenerateNewSecret())
		// r.Get("/fim", ar.FIM())

		r.Post("keys", ar.UploadJWTKeys())
		r.Get("keys", ar.GetJWTKeys())
	})

	r.Route("/invites", func(r chi.Router) {
		r.Use(ar.token())
		r.Get("/", ar.FetchInvites())
		r.Post("/", ar.AddInvite())
		r.Get("/{id:[a-zA-Z0-9]+}", ar.GetInviteByID())
		r.Delete("/{id:[a-zA-Z0-9]+}", ar.ArchiveInviteByID())
	})
	ar.router = r
}

func (ar *Router) token() func(next http.Handler) http.Handler {
	return wm.Token(model.TokenTypeManagement, ar.server.Services().Token, ar.server.Storages().Token, ar.LocalizedRouter)
}
