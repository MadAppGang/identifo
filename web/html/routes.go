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

	// If serve new web static files than just set web handlers and return
	if ar.staticFilesStorageSettings.ServeNewWeb {
		appHandler := ar.staticFilesStorage.WebHandlers()
		ar.Router.PathPrefix(`/`).Handler(appHandler.AppHandler).Methods("GET")

		return
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

	ar.Router.HandleFunc(`/{register:register/?}`, ar.HTMLFileHandler(model.StaticPagesNames.Registration)).Methods("GET")
	ar.Router.HandleFunc(`/password/{forgot:forgot/?}`, ar.HTMLFileHandler(model.StaticPagesNames.ForgotPassword)).Methods("GET")
	ar.Router.HandleFunc(`/password/forgot/{success:success/?}`, ar.HTMLFileHandler(model.StaticPagesNames.ForgotPasswordSuccess)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{error:error/?}`, ar.HTMLFileHandler(model.StaticPagesNames.TokenError)).Methods("GET")
	ar.Router.HandleFunc(`/password/reset/{success:success/?}`, ar.HTMLFileHandler(model.StaticPagesNames.ResetPasswordSuccess)).Methods("GET")
	ar.Router.HandleFunc(`/tfa/disable/{success:success/?}`, ar.HTMLFileHandler(model.StaticPagesNames.DisableTFASuccess)).Methods("GET")
	ar.Router.HandleFunc(`/tfa/reset/{success:success/?}`, ar.HTMLFileHandler(model.StaticPagesNames.ResetTFASuccess)).Methods("GET")
	ar.Router.HandleFunc(`/{misconfiguration:misconfiguration/?}`, ar.HTMLFileHandler(model.StaticPagesNames.Misconfiguration)).Methods("GET")

	assetHandlers := ar.staticFilesStorage.AssetHandlers()
	ar.Router.PathPrefix(`/{css:css/?}`).Handler(assetHandlers.StylesHandler).Methods("GET")
	ar.Router.PathPrefix(`/{js:js/?}`).Handler(assetHandlers.ScriptsHandler).Methods("GET")
	ar.Router.PathPrefix(`/{img:img/?}`).Handler(assetHandlers.ImagesHandler).Methods("GET")
	ar.Router.PathPrefix(`/{fonts:fonts/?}`).Handler(assetHandlers.FontsHandler).Methods("GET")
}
