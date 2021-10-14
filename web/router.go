package web

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/spa"
	"github.com/rs/cors"
)

const (
	adminpanelPath    = "/adminpanel"
	adminpanelAPIPath = "/admin"
	apiPath           = "/api"
	loginAppPath      = "/web"
)

// RouterSetting contains settings for root http router.
type RouterSetting struct {
	Server           model.Server
	Logger           *log.Logger
	ServeAdminPanel  bool
	HostName         string
	LoggerSettings   model.LoggerSettings
	AppOriginChecker model.OriginChecker
	APICors          *cors.Cors
	RestartChan      chan<- bool
}

// NewRouter creates and inits root http router.
func NewRouter(settings RouterSetting) (model.Router, error) {
	r := Router{}
	var err error
	authorizer := authorization.NewAuthorizer()

	// API router setup
	apiCorsSettings := model.DefaultCors
	if settings.AppOriginChecker != nil {
		apiCorsSettings.AllowOriginRequestFunc = settings.AppOriginChecker.CheckOrigin
	}
	apiCors := cors.New(apiCorsSettings)

	apiSettings := api.RouterSettings{
		Server:         settings.Server,
		Logger:         settings.Logger,
		LoggerSettings: settings.LoggerSettings,
		Authorizer:     authorizer,
		Host:           settings.HostName,
		Prefix:         apiPath,
		LoginWith:      settings.Server.Settings().Login.LoginWith,
		TFAType:        settings.Server.Settings().Login.TFAType,
		Cors:           apiCors,
	}

	apiRouter, err := api.NewRouter(apiSettings)
	if err != nil {
		return nil, err
	}
	r.APIRouter = apiRouter

	// Web login app setup
	loginAppSettings := spa.SPASettings{
		Root:       "/",
		FileSystem: http.FS(settings.Server.Storages().LoginAppFS),
	}
	r.LoginAppRouter, err = spa.NewRouter(loginAppSettings, settings.Logger)
	if err != nil {
		return nil, err
	}

	// Admin panel
	if settings.ServeAdminPanel {
		routerSettings := admin.RouterSettings{
			Server:  settings.Server,
			Logger:  settings.Logger,
			Host:    settings.HostName,
			Prefix:  adminpanelAPIPath,
			Restart: settings.RestartChan,
			OriginUpdate: func() error {
				return settings.AppOriginChecker.Update()
			},
		}

		// init admin panel api router
		r.AdminRouter, err = admin.NewRouter(routerSettings)
		if err != nil {
			return nil, err
		}
		// init admin panel web app
		adminPanelAppSettings := spa.SPASettings{
			Root:       "/",
			FileSystem: http.FS(settings.Server.Storages().AdminPanelFS),
		}
		r.AdminPanelRouter, err = spa.NewRouter(adminPanelAppSettings, settings.Logger)
		if err != nil {
			return nil, err
		}
		r.AdminPanelRouterPath = adminpanelPath
	}

	r.APIRouterPath = apiPath
	r.LoginAppRouterPath = loginAppPath

	r.setupRoutes()
	return &r, nil
}

// Router is a root router to handle REST API, web, and admin requests.
type Router struct {
	APIRouter        model.Router
	LoginAppRouter   model.Router
	AdminRouter      model.Router
	AdminPanelRouter model.Router
	RootRouter       *http.ServeMux

	APIRouterPath        string
	LoginAppRouterPath   string
	AdminRouterPath      string
	AdminPanelRouterPath string
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.RootRouter.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter = http.NewServeMux()
	ar.RootRouter.Handle("/", ar.APIRouter)
	ar.RootRouter.Handle(ar.LoginAppRouterPath+"/", http.StripPrefix(ar.LoginAppRouterPath, ar.LoginAppRouter))
	if ar.AdminRouter != nil && ar.AdminPanelRouter != nil {
		ar.RootRouter.Handle(ar.AdminRouterPath+"/", http.StripPrefix(ar.AdminRouterPath, ar.AdminRouter))
		ar.RootRouter.Handle(ar.AdminPanelRouterPath+"/", http.StripPrefix(ar.AdminPanelRouterPath, ar.AdminPanelRouter))
	}
}
