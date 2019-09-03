package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	configStoreEtcd "github.com/madappgang/identifo/configuration/storage/etcd"
	configStoreFile "github.com/madappgang/identifo/configuration/storage/file"
	configStoreS3 "github.com/madappgang/identifo/configuration/storage/s3"
	"github.com/madappgang/identifo/external_services/mail/mailgun"
	emailMock "github.com/madappgang/identifo/external_services/mail/mock"
	"github.com/madappgang/identifo/external_services/mail/ses"
	smsMock "github.com/madappgang/identifo/external_services/sms/mock"
	"github.com/madappgang/identifo/external_services/sms/twilio"
	ijwt "github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	dynamodb "github.com/madappgang/identifo/sessions/dynamodb"
	mem "github.com/madappgang/identifo/sessions/mem"
	redis "github.com/madappgang/identifo/sessions/redis"
	staticStoreLocal "github.com/madappgang/identifo/static/storage/local"
	staticStoreS3 "github.com/madappgang/identifo/static/storage/s3"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
	"gopkg.in/yaml.v2"
)

const serverConfigPathEnvName = "SERVER_CONFIG_PATH"

const (
	defaultAdminLogin    = "admin@admin.com"
	defaultAdminPassword = "password"
)

// ServerSettings are server settings.
var ServerSettings model.ServerSettings

func init() {
	loadServerConfigurationFromFile(&ServerSettings)
}

// loadServerConfigurationFromFile loads configuration from the yaml file and writes it to out variable.
func loadServerConfigurationFromFile(out *model.ServerSettings) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get server configuration file:", err)
	}

	// Iterate through possible config paths until we find the valid one.
	configPaths := []string{
		os.Getenv(serverConfigPathEnvName),
		"./server-config.yaml",
		"../../server/server-config.yaml",
	}

	var configFile []byte

	for _, p := range configPaths {
		if p == "" {
			continue
		}
		configFile, err = ioutil.ReadFile(filepath.Join(dir, p))
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatalln("Cannot read server configuration file:", err)
	}

	if err = yaml.Unmarshal(configFile, out); err != nil {
		log.Fatalln("Cannot unmarshal server configuration file:", err)
	}

	if len(os.Getenv(out.AdminAccount.LoginEnvName)) == 0 {
		if err := os.Setenv(out.AdminAccount.LoginEnvName, defaultAdminLogin); err != nil {
			log.Fatalf("Could not set default %s: %s\n", out.AdminAccount.LoginEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", out.AdminAccount.LoginEnvName)
	}
	if len(os.Getenv(out.AdminAccount.PasswordEnvName)) == 0 {
		if err := os.Setenv(out.AdminAccount.PasswordEnvName, defaultAdminPassword); err != nil {
			log.Fatalf("Could not set default %s: %s\n", out.AdminAccount.PasswordEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", out.AdminAccount.PasswordEnvName)
	}

	if err := out.Validate(); err != nil {
		log.Fatalln(err)
	}

	if err = os.Setenv(serverConfigPathEnvName, out.StaticFilesStorage.ServerConfigPath); err != nil {
		log.Println("Could not set server config path env variable. Strange yet not critical. Error:", err)
	}
}

// NewServer creates backend service.
func NewServer(settings model.ServerSettings, db DatabaseComposer, configurationStorage model.ConfigurationStorage, options ...func(*Server) error) (model.Server, error) {
	var err error
	if configurationStorage == nil {
		configurationStorage, err = InitConfigurationStorage(settings.ConfigurationStorage, settings.StaticFilesStorage.ServerConfigPath)
		if err != nil {
			return nil, err
		}
	}

	appStorage, userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, err := db.Compose()
	if err != nil {
		return nil, err
	}

	tokenService, err := initTokenService(settings.General, configurationStorage, tokenStorage, appStorage, userStorage)
	if err != nil {
		return nil, err
	}

	staticFilesStorage, err := initStaticFilesStorage(settings.StaticFilesStorage)
	if err != nil {
		return nil, err
	}

	s := Server{
		appStorage:              appStorage,
		userStorage:             userStorage,
		tokenStorage:            tokenStorage,
		tokenBlacklist:          tokenBlacklist,
		verificationCodeStorage: verificationCodeStorage,
		configurationStorage:    configurationStorage,
		staticFilesStorage:      staticFilesStorage,
	}

	sessionStorage, err := initSessionStorage(settings)
	if err != nil {
		return nil, err
	}
	sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, sessionStorage)

	ms, err := initEmailService(settings.ExternalServices.EmailService, staticFilesStorage)
	if err != nil {
		return nil, err
	}

	sms, err := initSMSService(settings.ExternalServices.SMSService)
	if err != nil {
		return nil, err
	}

	// env variable can rewrite host option
	hostName := os.Getenv("HOST_NAME")
	if len(hostName) == 0 {
		hostName = settings.General.Host
	}

	routerSettings := web.RouterSetting{
		AppStorage:              appStorage,
		UserStorage:             userStorage,
		TokenStorage:            tokenStorage,
		VerificationCodeStorage: verificationCodeStorage,
		TokenService:            tokenService,
		TokenBlacklist:          tokenBlacklist,
		SessionService:          sessionService,
		SessionStorage:          sessionStorage,
		ConfigurationStorage:    configurationStorage,
		StaticFilesStorage:      staticFilesStorage,
		ServeAdminPanel:         settings.StaticFilesStorage.ServeAdminPanel,
		SMSService:              sms,
		EmailService:            ms,
		WebRouterSettings: []func(*html.Router) error{
			html.HostOption(hostName),
		},
		APIRouterSettings: []func(*api.Router) error{
			api.HostOption(hostName),
			api.SupportedLoginWaysOption(settings.Login.LoginWith),
			api.TFATypeOption(settings.Login.TFAType),
		},
		AdminRouterSettings: []func(*admin.Router) error{
			admin.HostOption(hostName),
			admin.ServerConfigPathOption(settings.StaticFilesStorage.ServerConfigPath),
			admin.ServerSettingsOption(&settings),
		},
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
}

// Router returns server's main router.
func (s *Server) Router() model.Router {
	return s.MainRouter
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

// InitConfigurationStorage initializes configuration storage.
func InitConfigurationStorage(settings model.ConfigurationStorageSettings, serverConfigPath string) (model.ConfigurationStorage, error) {
	switch settings.Type {
	case model.ConfigurationStorageTypeEtcd:
		return configStoreEtcd.NewConfigurationStorage(settings, serverConfigPath)
	case model.ConfigurationStorageTypeS3:
		return configStoreS3.NewConfigurationStorage(settings)
	case model.ConfigurationStorageTypeFile:
		return configStoreFile.NewConfigurationStorage(settings)
	default:
		return nil, model.ErrorNotImplemented
	}
}

func initTokenService(generalSettings model.GeneralServerSettings, configStorage model.ConfigurationStorage, tokenStorage model.TokenStorage, appStorage model.AppStorage, userStorage model.UserStorage) (jwtService.TokenService, error) {
	tokenServiceAlg, ok := ijwt.StrToTokenSignAlg[generalSettings.Algorithm]
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

func initSessionStorage(settings model.ServerSettings) (model.SessionStorage, error) {
	switch settings.SessionStorage.Type {
	case model.SessionStorageRedis:
		return redis.NewSessionStorage(settings.SessionStorage)
	case model.SessionStorageMem:
		return mem.NewSessionStorage()
	case model.SessionStorageDynamoDB:
		return dynamodb.NewSessionStorage(settings.SessionStorage)
	default:
		return nil, model.ErrorNotImplemented
	}
}

func initStaticFilesStorage(settings model.StaticFilesStorageSettings) (model.StaticFilesStorage, error) {
	switch settings.Type {
	case model.StaticFilesStorageTypeLocal:
		return staticStoreLocal.NewStaticFilesStorage(settings)
	case model.StaticFilesStorageTypeS3:
		return staticStoreS3.NewStaticFilesStorage(settings)
	default:
		return nil, model.ErrorNotImplemented
	}
}

func initSMSService(settings model.SMSServiceSettings) (model.SMSService, error) {
	switch settings.Type {
	case model.SMSServiceTwilio:
		return twilio.NewSMSService(settings)
	case model.SMSServiceMock:
		return smsMock.NewSMSService()
	default:
		return nil, model.ErrorNotImplemented
	}
}

func initEmailService(ess model.EmailServiceSettings, sfs model.StaticFilesStorage) (model.EmailService, error) {
	tpltr, err := model.NewEmailTemplater(sfs)
	if err != nil {
		return nil, err
	}
	if tpltr == nil {
		return nil, errors.New("Email templater holds nil value")
	}

	switch ess.Type {
	case model.EmailServiceMailgun:
		return mailgun.NewEmailService(ess, tpltr), nil
	case model.EmailServiceAWS:
		return ses.NewEmailService(ess, tpltr)
	case model.EmailServiceMock:
		return emailMock.NewEmailService(), nil
	default:
		return nil, model.ErrorNotImplemented
	}
}

// ImportApps imports apps from file.
func (s *Server) ImportApps(filename string) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return s.AppStorage().ImportJSON(data)
}

// ImportUsers imports users from file.
func (s *Server) ImportUsers(filename string) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return s.UserStorage().ImportJSON(data)
}

func dataFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
