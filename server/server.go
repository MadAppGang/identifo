package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/madappgang/identifo/configuration/etcd/storage"
	configStoreMock "github.com/madappgang/identifo/configuration/mock/storage"
	"github.com/madappgang/identifo/external_services/mail/mailgun"
	"github.com/madappgang/identifo/external_services/mail/ses"
	smsMock "github.com/madappgang/identifo/external_services/sms/mock"
	"github.com/madappgang/identifo/external_services/sms/twilio"
	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	mem "github.com/madappgang/identifo/sessions/mem"
	redis "github.com/madappgang/identifo/sessions/redis"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
	"gopkg.in/yaml.v2"
)

const serverConfigPathEnv = "SERVER_CONFIG_PATH"

// ServerSettings are server settings.
var ServerSettings model.ServerSettings

func init() {
	LoadServerConfiguration(&ServerSettings)
}

// LoadServerConfiguration loads configuration from the yaml file and writes it to out variable.
func LoadServerConfiguration(out *model.ServerSettings) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get server configuration file:", err)
	}

	configPaths := []string{
		os.Getenv(serverConfigPathEnv),
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

	if len(out.AdminAccount.LoginEnvName) == 0 {
		log.Fatalln("Admin login env variable name not specified")
	}
	if len(out.AdminAccount.PasswordEnvName) == 0 {
		log.Fatalln("Admin password env variable name not specified")
	}

	if len(os.Getenv(out.AdminAccount.LoginEnvName)) == 0 {
		log.Fatalln("Admin login env variable not set")
	}
	if len(os.Getenv(out.AdminAccount.PasswordEnvName)) == 0 {
		log.Fatalln("Admin password env variable not set")
	}

	if err = os.Setenv(serverConfigPathEnv, out.ServerConfigPath); err != nil {
		log.Println("Could not set server config path env variable. Strange yet not critical. Error:", err)
	}
}

// DatabaseComposer inits database stack.
type DatabaseComposer interface {
	Compose() (
		model.AppStorage,
		model.UserStorage,
		model.TokenStorage,
		model.VerificationCodeStorage,
		jwtService.TokenService,
		error,
	)
}

// PartialDatabaseComposer can init services backed with different databases.
type PartialDatabaseComposer interface {
	AppStorageComposer() func() (model.AppStorage, error)
	UserStorageComposer() func() (model.UserStorage, error)
	TokenStorageComposer() func() (model.TokenStorage, error)
	VerificationCodeStorageComposer() func() (model.VerificationCodeStorage, error)
}

// Composer is a service composer which is agnostic to particular database implementations.
type Composer struct {
	settings                   model.ServerSettings
	newAppStorage              func() (model.AppStorage, error)
	newUserStorage             func() (model.UserStorage, error)
	newTokenStorage            func() (model.TokenStorage, error)
	newVerificationCodeStorage func() (model.VerificationCodeStorage, error)
}

// Compose composes all services.
func (c *Composer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {
	appStorage, err := c.newAppStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	userStorage, err := c.newUserStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := c.newTokenStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	verificationCodeStorage, err := c.newVerificationCodeStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenSignAlg[c.settings.Algorithm]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("Unknown token service algorithm %s", c.settings.Algorithm)
	}

	tokenService, err := jwtService.NewJWTokenService(
		path.Join(c.settings.PEMFolderPath, c.settings.PrivateKey),
		path.Join(c.settings.PEMFolderPath, c.settings.PublicKey),
		c.settings.Issuer,
		tokenServiceAlg,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, verificationCodeStorage, tokenService, nil
}

// NewComposer returns new database composer based on passed server settings.
func NewComposer(settings model.ServerSettings, partialComposers []PartialDatabaseComposer, options ...func(*Composer) error) (*Composer, error) {
	c := &Composer{settings: settings}

	for _, pc := range partialComposers {
		if pc.AppStorageComposer() != nil {
			c.newAppStorage = pc.AppStorageComposer()
		}
		if pc.UserStorageComposer() != nil {
			c.newUserStorage = pc.UserStorageComposer()
		}
		if pc.TokenStorageComposer() != nil {
			c.newTokenStorage = pc.TokenStorageComposer()
		}
		if pc.VerificationCodeStorageComposer() != nil {
			c.newVerificationCodeStorage = pc.VerificationCodeStorageComposer()
		}
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// NewServer creates backend service.
func NewServer(settings model.ServerSettings, db DatabaseComposer, options ...func(*Server) error) (model.Server, error) {
	appStorage, userStorage, tokenStorage, verificationCodeStorage, tokenService, err := db.Compose()
	if err != nil {
		return nil, err
	}

	configurationStorage, err := configurationStorage(settings.ConfigurationStorage)
	if err != nil {
		return nil, err
	}

	s := Server{appStorage: appStorage, userStorage: userStorage, configurationStorage: configurationStorage}

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
	MainRouter           *web.Router
	appStorage           model.AppStorage
	userStorage          model.UserStorage
	configurationStorage model.ConfigurationStorage
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

// ConfigurationStorage returns server's configuration storage.
func (s *Server) ConfigurationStorage() model.ConfigurationStorage {
	return s.configurationStorage
}

// ConfigurationStorageOption is an option to set server's configuration storage.
func ConfigurationStorageOption(configuratonStorage model.ConfigurationStorage) func(*Server) error {
	return func(s *Server) error {
		if configuratonStorage != nil {
			s.configurationStorage = configuratonStorage
		}
		return nil
	}
}

func configurationStorage(settings model.ConfigurationStorageSettings) (model.ConfigurationStorage, error) {
	switch settings.Type {
	case model.ConfigurationStorageTypeEtcd:
		return etcd.NewConfigurationStorage(settings)
	case model.ConfigurationStorageTypeMock:
		return configStoreMock.NewConfigurationStorage()
	default:
		return nil, model.ErrorNotImplemented
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
