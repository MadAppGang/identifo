package spa

import (
	"log/slog"
	"net/http"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
)

func NewRouter(setting SPASettings, middlewares []negroni.Handler) (model.Router, error) {
	ar := Router{
		FS: setting.FileSystem,
	}

	ar.Logger = logging.NewLogger(
		setting.LoggerSettings.Format,
		setting.LoggerSettings.SPA.Level,
	).With(logging.FieldComponent, setting.Name)

	ar.Middleware = buildMiddleware(
		setting.Name,
		setting.LoggerSettings.DumpRequest,
		setting.LoggerSettings.Format,
		setting.LoggerSettings.MaxBodySize,
		setting.LoggerSettings.SPA,
		!setting.LoggerSettings.LogSensitiveData,
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
	format string,
	maxBodySize int,
	logParams model.LoggerParams,
	logSensitiveData bool,
	middlewares []negroni.Handler,
) *negroni.Negroni {
	lm := middleware.NegroniHTTPLogger(
		settingName,
		format,
		maxBodySize,
		logParams,
		model.HTTPLogDetailing(dumpRequest, logParams.HTTPDetailing),
		!logSensitiveData,
	)

	handlers := []negroni.Handler{
		lm,
		negroni.NewRecovery(),
	}

	return negroni.New(handlers...).With(middlewares...)
}
