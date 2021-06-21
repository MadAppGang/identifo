package server

import (
	"fmt"
	"net/http"
	"os"

	jwtService "github.com/madappgang/identifo/jwt/service"
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
func NewServer(storages model.ServerStorageCollection, services model.ServerThirdPartyServices, options ...func(*Server) error) (model.Server, error) {
	if storages.Config == nil {
		return nil, fmt.Errorf("unable create sever with empty config storage")
	}

	settings, err := storages.Config.LoadServerSettings(false)
	if err != nil {
		return nil, err
	}

	s := Server{
		appStorage:              storages.App,
		userStorage:             storages.User,
		tokenStorage:            storages.Token,
		tokenBlacklist:          storages.Blocklist,
		verificationCodeStorage: storages.Verification,
		configurationStorage:    storages.Config,
		staticFilesStorage:      storages.Static,
		settings:                settings,
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

	sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, storages.Session)

	routerSettings := web.RouterSetting{
		AppStorage:              storages.App,
		UserStorage:             storages.User,
		TokenStorage:            storages.Token,
		VerificationCodeStorage: storages.Verification,
		TokenService:            tokenService,
		TokenBlacklist:          storages.Blocklist,
		InviteStorage:           storages.Invite,
		SessionService:          sessionService,
		SessionStorage:          storages.Session,
		ConfigurationStorage:    storages.Config,
		StaticFilesStorage:      storages.Static,
		ServeAdminPanel:         settings.StaticFilesStorage.ServeAdminPanel,
		SMSService:              services.SMS,
		EmailService:            services.Email,
		WebRouterSettings: []func(*html.Router) error{
			html.HostOption(hostName),
			html.StaticFilesStorageSettings(&settings.StaticFilesStorage),
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
			admin.ServerConfigPathOption(settings.StaticFilesStorage.ServerConfigPath),
			admin.ServerSettingsOption(&settings),
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
	MainRouter              *web.Router
	appStorage              model.AppStorage
	userStorage             model.UserStorage
	configurationStorage    model.ConfigurationStorage
	tokenStorage            model.TokenStorage
	tokenBlacklist          model.TokenBlacklist
	staticFilesStorage      model.StaticFilesStorage
	verificationCodeStorage model.VerificationCodeStorage
	settings                model.ServerSettings
}

// Router returns server's main router.
func (s *Server) Router() model.Router {
	return s.MainRouter
}

func (s *Server) Settings() model.ServerSettings {
	return s.settings
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.MainRouter.ServeHTTP(w, r)
}

// AppStorage returns server's app storage.
func (s *Server) AppStorage() model.AppStorage {
	return s.appStorage
}

// UserStorage returns server's user storage.
func (s *Server) UserStorage() model.UserStorage {
	return s.userStorage
}

// TokenStorage returns server's token storage.
func (s *Server) TokenStorage() model.TokenStorage {
	return s.tokenStorage
}

// TokenBlacklist returns server's token blacklist.
func (s *Server) TokenBlacklist() model.TokenBlacklist {
	return s.tokenBlacklist
}

// VerificationCodeStorage returns server's token storage.
func (s *Server) VerificationCodeStorage() model.VerificationCodeStorage {
	return s.verificationCodeStorage
}

// ConfigurationStorage returns server's configuration storage.
func (s *Server) ConfigurationStorage() model.ConfigurationStorage {
	return s.configurationStorage
}

// StaticFilesStorage returns server's static files storage.
func (s *Server) StaticFilesStorage() model.StaticFilesStorage {
	return s.staticFilesStorage
}

// Close closes all database connections.
func (s *Server) Close() {
	s.AppStorage().Close()
	s.UserStorage().Close()
	s.TokenStorage().Close()
	s.TokenBlacklist().Close()
	s.VerificationCodeStorage().Close()
	s.StaticFilesStorage().Close()
}

func initTokenService(generalSettings model.GeneralServerSettings, configStorage model.ConfigurationStorage, tokenStorage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage) (jwtService.TokenService, error) {
	tokenServiceAlg, ok := model.StrToTokenSignAlg[generalSettings.Algorithm]
	if !ok {
		return nil, fmt.Errorf("Unknown token service algorithm %s", generalSettings.Algorithm)
	}

	keys, err := configStorage.LoadKeys(tokenServiceAlg)
	if err != nil {
		return nil, err
	}

	tokenService, err := jwtService.NewJWTokenService(
		keys,
		generalSettings.Issuer,
		tokenStorage,
		appStorage,
		userStorage,
	)
	return tokenService, err
}
