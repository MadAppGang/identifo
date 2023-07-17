package spa

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madappgang/identifo/v2/model"
	mw "github.com/madappgang/identifo/v2/web/middleware"
)

func NewRouter(setting SPASettings, logger *log.Logger) (model.Router, error) {
	ar := Router{
		FS: setting.FileSystem,
	}

	// Setup logger to stdout.
	if logger == nil {
		ar.Logger = log.New(os.Stdout, fmt.Sprintf("[ %s ]: ", setting.Name), log.Ldate|log.Ltime|log.Lshortfile)
	}

	// ar.Middleware.UseHandler(NewSPAHandlerFunc(setting))

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	if setting.NewCacheDisabled {
		r.Use(mw.NewCacheDisable)
	}
	r.NotFound(NewSPAHandlerFunc(setting))
	ar.mux = r

	return &ar, nil
}

// login app router
type Router struct {
	Logger *log.Logger
	mux    *chi.Mux
	FS     http.FileSystem
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.mux.ServeHTTP(w, r)
}
