package web

import (
	"io/fs"
	"net/http"
	"net/url"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/web/admin"
	"github.com/madappgang/identifo/v2/web/api"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/management"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/madappgang/identifo/v2/web/spa"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const (
	adminpanelPath    = "/adminpanel"
	managementPath    = "/management"
	adminpanelAPIPath = "/admin"
	apiPath           = "/api"
)

// RouterSetting contains settings for root http router.
type RouterSetting struct {
	Server           model.Server
	ServeAdminPanel  bool
	Host             *url.URL
	LoggerSettings   model.LoggerSettings
	AppOriginChecker model.OriginChecker
	APICors          *cors.Cors
	RestartChan      chan<- bool
	Locale           string
}

// NewRootRouter creates and inits root http router.
func NewRootRouter(settings RouterSetting) (model.Router, error) {
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
		LoggerSettings:   settings.LoggerSettings,
		Authorizer:       authorizer,
		Host:             settings.Host,
		LoginWith:        settings.Server.Settings().Login.LoginWith,
		TFAType:          settings.Server.Settings().Login.TFAType,
		TFAResendTimeout: settings.Server.Settings().Login.TFAResendTimeout,
		Cors:             apiCors,
		Locale:           settings.Locale,
	}

	apiRouter, err := api.NewRouter(apiSettings)
	if err != nil {
		return nil, err
	}
	r.APIRouter = apiRouter

	managementRouter, err := management.NewRouter(management.RouterSettings{
		Server:             settings.Server,
		LoggerSettings:     settings.LoggerSettings,
		Storage:            settings.Server.Storages().ManagementKey,
		Locale:             settings.Locale,
		SupportedLoginWays: settings.Server.Settings().Login.LoginWith,
	})
	if err != nil {
		return nil, err
	}
	r.ManagementRouter = managementRouter

	if settings.Server.Settings().LoginWebApp.Type == model.FileStorageTypeNone {
		r.LoginAppRouter = nil
	} else {
		// Web login app setup
		loginAppSettings := spa.SPASettings{
			Name:           "LOGIN_APP",
			Root:           "/",
			FileSystem:     http.FS(settings.Server.Storages().LoginAppFS),
			LoggerSettings: settings.LoggerSettings,
		}
		r.LoginAppRouter, err = spa.NewRouter(
			loginAppSettings,
			[]negroni.Handler{middleware.NewCacheDisable()})
		if err != nil {
			return nil, err
		}
	}

	// Admin panel
	if settings.ServeAdminPanel {
		routerSettings := admin.RouterSettings{
			Server:         settings.Server,
			LoggerSettings: settings.LoggerSettings,
			Host:           settings.Host,
			Prefix:         adminpanelAPIPath,
			Restart:        settings.RestartChan,
		}

		if settings.AppOriginChecker != nil {
			checker := settings.AppOriginChecker // keep reference to origin checker, not settings
			routerSettings.OriginUpdate = func() error {
				return checker.Update()
			}
			r.UpdateCORS = func() {
				checker.Update()
			}
		}

		// init admin panel api router
		r.AdminRouter, err = admin.NewRouter(routerSettings)
		if err != nil {
			return nil, err
		}

		// init admin panel web app
		adminPanelAppSettings := spa.SPASettings{
			Name:           "ADMIN_PANEL",
			Root:           "/",
			FileSystem:     http.FS(fsWithConfig(settings.Server.Storages().AdminPanelFS)),
			LoggerSettings: settings.LoggerSettings,
		}
		r.AdminPanelRouter, err = spa.NewRouter(adminPanelAppSettings, nil)
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
	ManagementRouter model.Router
	LoginAppRouter   model.Router
	AdminRouter      model.Router
	AdminPanelRouter model.Router
	RootRouter       *http.ServeMux
	UpdateCORS       func()
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.RootRouter.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter = http.NewServeMux()
	ar.RootRouter.Handle("/", ar.APIRouter)
	if ar.ManagementRouter != nil {
		ar.RootRouter.Handle(managementPath+"/", http.StripPrefix(managementPath, ar.ManagementRouter))
	}
	if ar.LoginAppRouter != nil {
		ar.RootRouter.Handle(model.DefaultLoginWebAppSettings.LoginURL+"/", http.StripPrefix(model.DefaultLoginWebAppSettings.LoginURL, ar.LoginAppRouter))
	}
	if ar.AdminRouter != nil && ar.AdminPanelRouter != nil {
		ar.RootRouter.Handle(adminpanelAPIPath+"/", http.StripPrefix(adminpanelAPIPath, ar.AdminRouter))
		ar.RootRouter.Handle(adminpanelPath+"/", http.StripPrefix(adminpanelPath, ar.AdminPanelRouter))
	}
}
