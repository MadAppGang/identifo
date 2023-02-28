package management

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
)

// setup all routes
func (ar *Router) initRoutes() {
	if ar.router == nil {
		panic("Empty API router")
	}

	baseMiddleware := negroni.New(
		middleware.NewNegroniLogger("MANAGEMENT"),
		negroni.NewRecovery(),
		ar.RemoveTrailingSlash(),
	)

	if ar.loggerSettings.DumpRequest {
		baseMiddleware.Use(ar.DumpRequest())
	}

	ar.router.Handle("/test", with(baseMiddleware, negroni.WrapFunc(ar.test))).Methods(http.MethodGet)
	ar.router.Handle("/invite_token", with(baseMiddleware, negroni.WrapFunc(ar.getInviteToken))).Methods(http.MethodPost)
	ar.router.Handle("/reset_password_token", with(baseMiddleware, negroni.WrapFunc(ar.getResetPasswordToken))).Methods(http.MethodPost)
}

func (ar *Router) RemoveTrailingSlash() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(rw, r)
	}
}

// DumpRequest logs the request.
func (ar *Router) DumpRequest() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			ar.logger.Println("Error dumping request:", err)
		}
		ar.logger.Printf("Request: %s\n", string(dump))
		next(rw, r)
	}
}

func with(n *negroni.Negroni, handlers ...negroni.Handler) *negroni.Negroni {
	existing := n.Handlers()
	h := []negroni.Handler{}
	h = append(h, existing...)
	h = append(h, handlers...)
	return negroni.New(h...)
}
