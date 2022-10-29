package server

import (
	"fmt"
	"net/http"
	"os"

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

	// env variable can rewrite host option
	hostName := os.Getenv("IDENTIFO_HOST_NAME")
	if len(hostName) == 0 {
		hostName = settings.General.Host
	}

	originChecker, err := middleware.NewAppOriginChecker(storages.App)
	if err != nil {
		return nil, err
	}

	// sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, storages.Session)
	routerSettings := web.RouterSetting{
		Server:           &s,
		ServeAdminPanel:  settings.AdminPanel.Enabled,
		HostName:         hostName,
		AppOriginChecker: originChecker,
		RestartChan:      restartChan,
		LoggerSettings:   settings.Logger,
	}

	r, err := web.NewRouter(routerSettings)
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
	s.storages.App.Close()
	s.storages.User.Close()
	s.storages.Token.Close()
	s.storages.Blocklist.Close()
	s.storages.Invite.Close()
	s.storages.Verification.Close()
	s.storages.Session.Close()
}

func (s *Server) Errors() []error {
	return s.errs
}
