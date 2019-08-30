package html

import (
	"github.com/madappgang/identifo/model"
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

	ar.Router.Path(`/tfa/{reset:reset/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetTFA()),
	)).Methods("POST")

	ar.Router.Path(`/tfa/{reset:reset/?}`).Handler(negroni.New(
		ar.ResetTokenMiddleware(),
		negroni.WrapFunc(ar.ResetTFAHandler()),
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

	// ar.Router.HandleFunc(`/{register:register/?}`, ar.HTMLFileHandler(model.RegistrationStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/password/{forgot:forgot/?}`, ar.HTMLFileHandler(model.ForgotPasswordStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/password/forgot/{success:success/?}`, ar.HTMLFileHandler(model.ForgotPasswordSuccessStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{error:error/?}`, ar.HTMLFileHandler(model.TokenErrorStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{success:success/?}`, ar.HTMLFileHandler(model.ResetPasswordSuccessStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/tfa/disable/{success:success/?}`, ar.HTMLFileHandler(model.DisableTFASuccessStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/tfa/reset/{success:success/?}`, ar.HTMLFileHandler(model.ResetTFASuccessStaticPageName)).Methods("GET")
	ar.Router.HandleFunc(`/{misconfiguration:misconfiguration/?}`, ar.HTMLFileHandler(model.MisconfigurationStaticPageName)).Methods("GET")

	stylesHandler := ar.staticFilesStorage.StylesHandler()
	scriptsHandler := ar.staticFilesStorage.ScriptsHandler()
	imagesHandler := ar.staticFilesStorage.ImagesHandler()
	fontsHandler := ar.staticFilesStorage.FontsHandler()

	// Setup routes for static files.
	ar.Router.PathPrefix(`/{css:css/?}`).Handler(stylesHandler).Methods("GET")
	ar.Router.PathPrefix(`/{js:js/?}`).Handler(scriptsHandler).Methods("GET")
	ar.Router.PathPrefix(`/{img:img/?}`).Handler(imagesHandler).Methods("GET")
	ar.Router.PathPrefix(`/{fonts:fonts/?}`).Handler(fontsHandler).Methods("GET")
}
