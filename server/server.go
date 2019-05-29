package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/madappgang/identifo/external_services/mail/mailgun"
	"github.com/madappgang/identifo/external_services/mail/ses"
	"github.com/madappgang/identifo/external_services/sms/twilio"
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

const serverConfigPath = "../../server/server-config.yaml"

// ServerSettings are server settings.
var ServerSettings model.ServerSettings

func init() {
	LoadServerConfiguration(&ServerSettings)
}

// LoadServerConfiguration loads configuration from the yaml file and writes it to out variable.
func LoadServerConfiguration(out interface{}) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get server configuration file:", err)
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, serverConfigPath))
	var origErr = err
	if err != nil {
		yamlFile, err = ioutil.ReadFile(path.Base("./server-config.yaml"))
		if err != nil {
			log.Fatalln("Cannot read server configuration file:", origErr)
		}
	}

	if err = yaml.Unmarshal(yamlFile, out); err != nil {
		log.Fatalln("Cannot unmarshal configuration file:", err)
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

// NewServer creates backend service.
func NewServer(settings model.ServerSettings, db DatabaseComposer, options ...func(*Server) error) (model.Server, error) {
	appStorage, userStorage, tokenStorage, verificationCodeStorage, tokenService, err := db.Compose()
	if err != nil {
		return nil, err
	}
	s := Server{AppStrg: appStorage, UserStrg: userStorage}

	sessionStorage, err := sessionStorage(settings.SessionStorage)
	if err != nil {
		return nil, err
	}
	sessionService := model.NewSessionManager(settings.SessionDuration, sessionStorage)

	ms, err := mailService(settings.MailService, settings.EmailTemplateNames, settings.EmailTemplatesPath)
	if err != nil {
		return nil, err
	}

	sms, err := smsService(settings)
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
			admin.AccountConfigPathOption(settings.AccountConfigPath),
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

// Server is DynamoDB-backed server.
type Server struct {
	MainRouter *web.Router
	AppStrg    model.AppStorage
	UserStrg   model.UserStorage
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
	return s.AppStrg
}

// UserStorage returns server's user storage.
func (s *Server) UserStorage() model.UserStorage {
	return s.UserStrg
}

func smsService(settings model.ServerSettings) (model.SMSService, error) {
	switch settings.SMSService {
	case model.SMSServiceTwilio:
		return twilio.NewSMSService(settings.Twilio.AccountSid, settings.Twilio.AuthToken, settings.Twilio.ServiceSid)
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

func sessionStorage(storageType model.SessionStorageType) (model.SessionStorage, error) {
	switch storageType {
	case model.SessionStorageRedis:
		return redis.NewSessionStorageFromEnv()
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
