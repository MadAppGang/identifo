package html

import (
	"net/http"

	"github.com/urfave/negroni"
)

//setup all routes
func (ar *Router) initRoutes() {
	//do nothing on empty router (or should panic?)
	if ar.Router == nil {
		return
	}

	ar.Router.Path("/password/{reset:reset\\/?}").Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetPassword()),
	)).Methods("POST")

	ar.Router.Path("/password/{reset:reset\\/?}").Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetPasswordHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ResetPassword)),
	)).Methods("GET")

	ar.Router.HandleFunc("/password/{forgot:forgot\\/?}", ar.SendResetToken()).Methods("POST")

	ar.Router.Path("/{login:login\\/?}").Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Login()),
	)).Methods("POST")

	ar.Router.Path("/{login:login\\/?}").Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.LoginHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Login)),
	)).Methods("GET")

	ar.Router.Path("/{register:register\\/?}").Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Register()),
	)).Methods("POST")

	ar.Router.Path("/{register:register\\/?}").Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.RegistrationHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Registration)),
	)).Methods("GET")

	ar.Router.HandleFunc("/token/{renew:renew\\/?}", ar.RenewToken(ar.StaticFilesPath.PagesPath, "/web-message.html")).Methods("GET")
	ar.Router.Path("/{logout:logout\\/?}").Handler(negroni.New(
		ar.AppID(),
		negroni.WrapFunc(ar.Logout()),
	)).Methods("GET")

	// ar.Router.HandleFunc("/{register:register\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Registration)).Methods("GET")
	ar.Router.HandleFunc("/password/{forgot:forgot\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ForgotPassword)).Methods("GET")
	ar.Router.HandleFunc("/password/forgot/{success:success\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ForgotPasswordSuccess)).Methods("GET")
	ar.Router.HandleFunc("/password/reset/{error:error\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.TokenError)).Methods("GET")
	ar.Router.HandleFunc("/password/reset/{success:success\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.ResetSuccess)).Methods("GET")
	ar.Router.HandleFunc("/{misconfiguration:misconfiguration\\/?}", ar.HTMLFileHandler(ar.StaticFilesPath.PagesPath, ar.StaticPages.Misconfiguration)).Methods("GET")

	stylesHandler := http.FileServer(http.Dir(ar.StaticFilesPath.StylesPath))
	scriptsHandler := http.FileServer(http.Dir(ar.StaticFilesPath.ScriptsPath))
	imagesHandler := http.FileServer(http.Dir(ar.StaticFilesPath.ImagesPath))

	//setup routes for static files
	ar.Router.PathPrefix("/{css:css\\/?}").Handler(http.StripPrefix("/css/", stylesHandler)).Methods("GET")
	ar.Router.PathPrefix("/{js:js\\/?}").Handler(http.StripPrefix("/js/", scriptsHandler)).Methods("GET")
	ar.Router.PathPrefix("/{img:img\\/?}").Handler(http.StripPrefix("/img/", imagesHandler)).Methods("GET")

}
