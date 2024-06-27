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
		FS: setting.FileSystem,
	}

	// Setup logger to stdout.
	if logger == nil {
		logger = logging.NewDefaultLogger().With(
			logging.FieldComponent, setting.Name)
	}

	ar.Logger = logger

	ar.Middleware = buildMiddleware(
		setting.Name,
		setting.LoggerSettings.DumpRequest,
		setting.LoggerSettings.SPA,
		middlewares,
	)

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

func buildMiddleware(
	settingName string,
	dumpRequest bool,
	logParams model.LoggerParams,
	middlewares []negroni.Handler,
) *negroni.Negroni {
	var handlers []negroni.Handler

	// set efficient log type
	logParams.Type = model.LogType(dumpRequest, logParams.Type)

	lm := middleware.HTTPLogger(
		settingName,
		logParams)
	handlers = append(handlers, lm)

	handlers = append(handlers, negroni.NewRecovery())

	return negroni.New(handlers...).With(middlewares...)
}
