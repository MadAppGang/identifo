package server

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// NewServer creates backend service.
func NewServer(storages model.ServerStorageCollection, services model.ServerServices, errs []error, restartChan chan<- bool) (model.Server, error) {
	if storages.Config == nil {
		return nil, fmt.Errorf("unable create sever with empty config storage")
	}

	// should be loaded and validates
	settings := storages.Config.LoadedSettings()
	if settings == nil {
		// no settings and no errors
		if len(errs) == 0 {
			return nil, fmt.Errorf("New Server could not be created, no settings loaded and no errors detected")
		} else {
			settings = &model.DefaultServerSettings
		}
	}

	s := Server{
		storages: storages,
		services: services,
		settings: *settings,
		errs:     errs, // keep the list of errors which keep server invalid but working
	}

	var originChecker *middleware.AppOriginChecker
	if len(errs) == 0 {
		// we have valid config loaded and we can do origin checker
		var err error
		originChecker, err = middleware.NewAppOriginChecker(storages.App)
		if err != nil {
			return nil, err
		}
	}

	// sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, storages.Session)
	routerSettings := web.RouterSetting{
		Server:           &s,
		ServeAdminPanel:  settings.AdminPanel.Enabled,
		AppOriginChecker: originChecker,
		RestartChan:      restartChan,
		LoggerSettings:   settings.Logger,
		Locale:           settings.General.Locale,
	}

	r, err := web.NewRootRouter(routerSettings)
	if err != nil {
		return nil, err
	}
	s.MainRouter = r.(*web.Router)
	return &s, nil
}

// Server is a server.
type Server struct {
	MainRouter *web.Router
	storages   model.ServerStorageCollection
	services   model.ServerServices
	settings   model.ServerSettings
	errs       []error
}

// Router returns server's main router.
func (s *Server) Router() model.Router {
	return s.MainRouter
}

func (s *Server) UpdateCORS() {
	if s.MainRouter.UpdateCORS != nil {
		s.MainRouter.UpdateCORS()
	}
}

func (s *Server) Settings() model.ServerSettings {
	return s.settings
}

func (s *Server) Services() model.ServerServices {
	return s.services
}

func (s *Server) Storages() model.ServerStorageCollection {
	return s.storages
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.MainRouter.ServeHTTP(w, r)
}

// Close closes all database connections.
func (s *Server) Close() {
	maybeClose := func(c interface{ Close() }) {
		if c != nil {
			c.Close()
		}
	}

	maybeClose(s.storages.App)
	// TODO: Implement close for user storages
	// maybeClose(s.storages.User)
	maybeClose(s.storages.Token)
	maybeClose(s.storages.Blocklist)
	maybeClose(s.storages.Invite)
	maybeClose(s.storages.Session)
}

func (s *Server) Errors() []error {
	return s.errs
}
