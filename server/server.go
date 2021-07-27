package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/utils/originchecker"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
	"github.com/rs/cors"
)

var defaultCors = model.CorsOptions{
	API: &cors.Options{AllowedHeaders: []string{"*", "x-identifo-clientid"}, AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}},
}

// NewServer creates backend service.
func NewServer(storages model.ServerStorageCollection, services model.ServerServices, options ...func(*Server) error) (model.Server, error) {
	if storages.Config == nil {
		return nil, fmt.Errorf("unable create sever with empty config storage")
	}

	settings, err := storages.Config.LoadServerSettings(false)
	if err != nil {
		return nil, err
	}

	s := Server{
		storages: storages,
		services: services,
		settings: settings,
	}

	// env variable can rewrite host option
	hostName := os.Getenv("HOST_NAME")
	if len(hostName) == 0 {
		hostName = settings.General.Host
	}

	originChecker := originchecker.NewOriginChecker()

	// validate, try to fetch apps
	apps, _, err := storages.App.FetchApps("", 0, 0)
	if err != nil {
		return nil, err
	}

	for _, a := range apps {
		originChecker.AddRawURLs(a.RedirectURLs)
	}

	// sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, storages.Session)

	routerSettings := web.RouterSetting{
		Server:          &s,
		ServeAdminPanel: settings.Static.ServeAdminPanel,
		WebRouterSettings: []func(*html.Router) error{
			html.HostOption(hostName),
			// html.StaticFilesStorageSettings(&settings.StaticFilesStorage),
			html.CorsOption(defaultCors),
		},
		APIRouterSettings: []func(*api.Router) error{
			api.HostOption(hostName),
			api.SupportedLoginWaysOption(settings.Login.LoginWith),
			api.TFATypeOption(settings.Login.TFAType),
			api.CorsOption(&defaultCors, originChecker),
		},
		AdminRouterSettings: []func(*admin.Router) error{
			admin.HostOption(hostName),
			// admin.ServerConfigPathOption(settings.StaticFilesStorage.ServerConfigPath),
			// admin.ServerSettingsOption(&settings),
			admin.CorsOption(defaultCors, originChecker),
		},
		LoggerSettings: settings.Logger,
	}

	r, err := web.NewRouter(routerSettings)
	if err != nil {
		return nil, err
	}
	s.MainRouter = r.(*web.Router)

	for _, option := range options {
		if err := option(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Server is a server.
type Server struct {
	MainRouter *web.Router
	storages   model.ServerStorageCollection
	services   model.ServerServices
	settings   model.ServerSettings
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
	s.storages.Blocklist.Close()
	s.storages.Token.Close()
	s.storages.Verification.Close()
	s.storages.Static.Close()
}
