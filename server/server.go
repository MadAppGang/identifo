package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	configStoreEtcd "github.com/madappgang/identifo/configuration/storage/etcd"
	configStoreFile "github.com/madappgang/identifo/configuration/storage/file"
	configStoreS3 "github.com/madappgang/identifo/configuration/storage/s3"
	"github.com/madappgang/identifo/external_services/mail/mailgun"
	emailMock "github.com/madappgang/identifo/external_services/mail/mock"
	"github.com/madappgang/identifo/external_services/mail/ses"
	smsMock "github.com/madappgang/identifo/external_services/sms/mock"
	"github.com/madappgang/identifo/external_services/sms/twilio"
	"github.com/madappgang/identifo/model"
	dynamodb "github.com/madappgang/identifo/sessions/dynamodb"
	mem "github.com/madappgang/identifo/sessions/mem"
	redis "github.com/madappgang/identifo/sessions/redis"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/adminpanel"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
	"gopkg.in/yaml.v2"
)

const (
	serverConfigPathEnvName = "SERVER_CONFIG_PATH"
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
		log.Fatalf("%s not set\n", out.AdminAccount.LoginEnvName)
	}
	if len(os.Getenv(out.AdminAccount.PasswordEnvName)) == 0 {
		log.Fatalf("%s not set\n", out.AdminAccount.PasswordEnvName)
	}

	if err := out.Validate(); err != nil {
		log.Fatalln(err)
	}

	if err = os.Setenv(serverConfigPathEnvName, out.ServerConfigPath); err != nil {
		log.Println("Could not set server config path env variable. Strange yet not critical. Error:", err)
	}
}

// NewServer creates backend service.
func NewServer(settings model.ServerSettings, db DatabaseComposer, configurationStorage model.ConfigurationStorage, options ...func(*Server) error) (model.Server, error) {
	appStorage, userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, tokenService, err := db.Compose()
	if err != nil {
		return nil, err
	}

	if configurationStorage == nil {
		configurationStorage, err = InitConfigurationStorage(settings.ConfigurationStorage, settings.ServerConfigPath)
		if err != nil {
			return nil, err
		}
	}

	s := Server{
		appStorage:              appStorage,
		userStorage:             userStorage,
		tokenStorage:            tokenStorage,
		tokenBlacklist:          tokenBlacklist,
		verificationCodeStorage: verificationCodeStorage,
		configurationStorage:    configurationStorage,
	}

	sessionStorage, err := sessionStorage(settings)
	if err != nil {
		return nil, err
	}
	sessionService := model.NewSessionManager(settings.SessionStorage.SessionDuration, sessionStorage)

	ms, err := mailService(settings.MailService, settings.EmailTemplateNames, settings.EmailTemplatesPath)
	if err != nil {
		return nil, err
	}

	sms, err := smsService(settings.SMSService)
	if err != nil {
		return nil, err
	}

	// env variable can rewrite host option
	hostName := os.Getenv("HOST_NAME")
	if len(hostName) == 0 {
		hostName = settings.Host
	}

	staticFiles := html.StaticFilesPath{
		StylesPath:  path.Join(settings.StaticFolderPath, "css"),
		ScriptsPath: path.Join(settings.StaticFolderPath, "js"),
		PagesPath:   settings.StaticFolderPath,
		ImagesPath:  path.Join(settings.StaticFolderPath, "img"),
		FontsPath:   path.Join(settings.StaticFolderPath, "fonts"),
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
		SMSService:              sms,
		EmailService:            ms,
		WebRouterSettings: []func(*html.Router) error{
			html.StaticPathOptions(staticFiles),
			html.HostOption(hostName),
		},
		APIRouterSettings: []func(*api.Router) error{
			api.HostOption(hostName),
			api.SupportedLoginWaysOption(settings.LoginWith),
			api.TFATypeOption(settings.TFAType),
		},
		AdminRouterSettings: []func(*admin.Router) error{
			admin.HostOption(hostName),
			admin.ServerConfigPathOption(settings.ServerConfigPath),
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

// Close closes all database connections.
func (s *Server) Close() {
	s.AppStorage().Close()
	s.UserStorage().Close()
	s.TokenStorage().Close()
	s.TokenBlacklist().Close()
	s.VerificationCodeStorage().Close()
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

// ServeAdminPanelOption is an option to serve admin panel right from the Identifo server.
func ServeAdminPanelOption() func(*Server) error {
	return func(s *Server) (err error) {
		s.MainRouter.AdminPanelRouter, err = adminpanel.NewRouter(adminpanel.BuildPathOption(ServerSettings.AdminPanelBuildPath))
		if err != nil {
			return
		}

		s.MainRouter.AdminPanelRouterPath = "/adminpanel"
		s.MainRouter.RootRouter.Handle(s.MainRouter.AdminPanelRouterPath+"/", http.StripPrefix(s.MainRouter.AdminPanelRouterPath, s.MainRouter.AdminPanelRouter))

		return nil
	}
}

func smsService(settings model.SMSServiceSettings) (model.SMSService, error) {
	switch settings.Type {
	case model.SMSServiceTwilio:
		return twilio.NewSMSService(settings)
	case model.SMSServiceMock:
		return smsMock.NewSMSService()
	default:
		return nil, model.ErrorNotImplemented
	}
}

func mailService(serviceType model.MailServiceType, templateNames model.EmailTemplateNames, templatesPath string) (model.EmailService, error) {
	tpltr, err := model.NewEmailTemplater(templateNames, templatesPath)
	if err != nil {
		return nil, err
	}
	if tpltr == nil {
		return nil, errors.New("Email templater holds nil value")
	}

	switch serviceType {
	case model.MailServiceMailgun:
		return mailgun.NewEmailServiceFromEnv(tpltr), nil
	case model.MailServiceAWS:
		return ses.NewEmailServiceFromEnv(tpltr)
	case model.MailServiceMock:
		return emailMock.NewEmailService(), nil
	default:
		return nil, model.ErrorNotImplemented
	}
}

func sessionStorage(settings model.ServerSettings) (model.SessionStorage, error) {
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
