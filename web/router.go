package web

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/web/admin"
	"github.com/madappgang/identifo/v2/web/api"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/madappgang/identifo/v2/web/spa"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const (
	adminpanelPath    = "/adminpanel"
	adminpanelAPIPath = "/admin"
	apiPath           = "/api"
	loginAppPath      = "/web"
	loginAppErrorPath = "/web/misconfiguration"
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
		Server:           settings.Server,
		Logger:           settings.Logger,
		LoggerSettings:   settings.LoggerSettings,
		Authorizer:       authorizer,
		Host:             settings.HostName,
		LoginAppPath:     loginAppPath,
		LoginWith:        settings.Server.Settings().Login.LoginWith,
		TFAType:          settings.Server.Settings().Login.TFAType,
		TFAResendTimeout: settings.Server.Settings().Login.TFAResendTimeout,
		Cors:             apiCors,
	}

	apiRouter, err := api.NewRouter(apiSettings)
	if err != nil {
		return nil, err
	}
	r.APIRouter = apiRouter

	if settings.Server.Settings().LoginWebApp.Type == model.FileStorageTypeNone {
		r.LoginAppRouter = nil
	} else {
		// Web login app setup
		loginAppSettings := spa.SPASettings{
			Name:       "LOGIN_APP",
			Root:       "/",
			FileSystem: http.FS(settings.Server.Storages().LoginAppFS),
		}
		r.LoginAppRouter, err = spa.NewRouter(loginAppSettings, []negroni.Handler{middleware.NewCacheDisable()}, settings.Logger)
		if err != nil {
			return nil, err
		}
	}

	// Admin panel
	if settings.ServeAdminPanel {
		routerSettings := admin.RouterSettings{
			Server:       settings.Server,
			Logger:       settings.Logger,
			Host:         settings.HostName,
			Prefix:       adminpanelAPIPath,
			Restart:      settings.RestartChan,
			LoginAppPath: loginAppPath,
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
			Name:       "ADMIN_PANEL",
			Root:       "/",
			FileSystem: http.FS(fsWithConfig(settings.Server.Storages().AdminPanelFS)),
		}
		r.AdminPanelRouter, err = spa.NewRouter(adminPanelAppSettings, nil, settings.Logger)
		if err != nil {
			return nil, err
		}
	}

	r.setupRoutes()
	return &r, nil
}

func fsWithConfig(fs fs.FS) fs.FS {
	files := map[string][]byte{
		"config.json": []byte(`{"apiUrl": "/admin"}`),
	}
	return mem.NewMapOverlayFS(fs, files)
}

// Router is a root router to handle REST API, web, and admin requests.
type Router struct {
	APIRouter        model.Router
	LoginAppRouter   model.Router
	AdminRouter      model.Router
	AdminPanelRouter model.Router
	RootRouter       *http.ServeMux
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.RootRouter.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter = http.NewServeMux()
	ar.RootRouter.Handle("/", ar.APIRouter)
	if ar.LoginAppRouter != nil {
		ar.RootRouter.Handle(loginAppPath+"/", http.StripPrefix(loginAppPath, ar.LoginAppRouter))
	}
	if ar.AdminRouter != nil && ar.AdminPanelRouter != nil {
		ar.RootRouter.Handle(adminpanelAPIPath+"/", http.StripPrefix(adminpanelAPIPath, ar.AdminRouter))
		ar.RootRouter.Handle(adminpanelPath+"/", http.StripPrefix(adminpanelPath, ar.AdminPanelRouter))
	}
}
