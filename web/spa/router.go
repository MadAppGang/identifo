package spa

import (
	"log/slog"
	"net/http"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
)

func NewRouter(setting SPASettings, middlewares []negroni.Handler, logger *slog.Logger) (model.Router, error) {
	ar := Router{
		Middleware: negroni.New(
			middleware.NewNegroniLogger(setting.Name),
			negroni.NewRecovery()).With(middlewares...),
		FS: setting.FileSystem,
	}

	// Setup logger to stdout.
	if logger == nil {
		ar.Logger = logging.NewDefaultLogger().With(
			logging.FieldComponent, setting.Name)
	}

	ar.Middleware.UseHandler(NewSPAHandlerFunc(setting))
	return &ar, nil
}

// login app router
type Router struct {
	Logger     *slog.Logger
	Middleware *negroni.Negroni
	FS         http.FileSystem
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.Middleware.ServeHTTP(w, r)
}
