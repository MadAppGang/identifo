package html

import (
	"net/http"

	"github.com/urfave/negroni"
)

// Setup all routes.
func (ar *Router) initRoutes() {
	if ar.Router == nil {
		panic("Empty HTML router")
	}

	ar.Router.Path(`/password/{reset:reset/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetPassword()),
	)).Methods("POST")

	ar.Router.Path(`/password/{reset:reset/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetPasswordHandler()),
	)).Methods("GET")

	ar.Router.Path(`/tfa/{disable:disable/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.DisableTFA()),
	)).Methods("POST")

	ar.Router.Path(`/tfa/{disable:disable/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.DisableTFAHandler()),
	)).Methods("GET")

	ar.Router.HandleFunc(`/password/{forgot:forgot/?}`, ar.SendResetToken()).Methods("POST")

	ar.Router.Path(`/{login:login/?}`).Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")

	ar.Router.Path(`/{login:login/?}`).Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.LoginHandler()),
	)).Methods("GET")

	ar.Router.Path(`/{register:register/?}`).Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Register()),
	)).Methods("POST")

	ar.Router.Path(`/{register:register/?}`).Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.RegistrationHandler()),
	)).Methods("GET")

	ar.Router.HandleFunc(`/token/{renew:renew/?}`, ar.RenewToken()).Methods("GET")
	ar.Router.Path(`/{logout:logout/?}`).Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Logout()),
	)).Methods("GET")

	// ar.Router.HandleFunc(`/{register:register/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Registration)).Methods("GET")
	ar.Router.HandleFunc(`/password/{forgot:forgot/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ForgotPassword)).Methods("GET")
	ar.Router.HandleFunc(`/password/forgot/{success:success/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ForgotPasswordSuccess)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{error:error/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.TokenError)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{success:success/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ResetPasswordSuccess)).Methods("GET")
	ar.Router.HandleFunc(`/tfa/disable/success/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.DisableTFASuccess)).Methods("GET")
	ar.Router.HandleFunc(`/{misconfiguration:misconfiguration/?}`, ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Misconfiguration)).Methods("GET")

	stylesHandler := http.FileServer(http.Dir(ar.StaticFilesPath.StylesPath))
	scriptsHandler := http.FileServer(http.Dir(ar.StaticFilesPath.ScriptsPath))
	imagesHandler := http.FileServer(http.Dir(ar.StaticFilesPath.ImagesPath))
	fontsHandler := http.FileServer(http.Dir(ar.StaticFilesPath.FontsPath))

	// Setup routes for static files.
	ar.Router.PathPrefix(`/{css:css/?}`).Handler(http.StripPrefix("/css/", stylesHandler)).Methods("GET")
	ar.Router.PathPrefix(`/{js:js/?}`).Handler(http.StripPrefix("/js/", scriptsHandler)).Methods("GET")
	ar.Router.PathPrefix(`/{img:img/?}`).Handler(http.StripPrefix("/img/", imagesHandler)).Methods("GET")
	ar.Router.PathPrefix(`/{fonts:img/?}`).Handler(http.StripPrefix("/fonts/", fontsHandler)).Methods("GET")
}
