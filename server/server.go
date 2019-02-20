package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/madappgang/identifo/mailgun"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/ses"
	mem "github.com/madappgang/identifo/sessions/mem"
	redis "github.com/madappgang/identifo/sessions/redis"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
)

// DefaultSettings are the default server settings.
var DefaultSettings = model.ServerSettings{
	StaticFolderPath:   "./static",
	PEMFolderPath:      "./pem",
	PrivateKey:         "private.pem",
	PublicKey:          "public.pem",
	EncryptionKeyPath:  "./encryptor/key.key",
	Algorithm:          model.TokenServiceAlgorithmAuto,
	Issuer:             "identifo",
	MailService:        model.MailServiceMailgun,
	SessionStorage:     model.SessionStorageMem,
	SessionDuration:    5 * time.Minute,
	Host:               "http://localhost:8080",
	ConfigPath:         "../../server/conf.yaml",
	EmailTemplatesPath: "./email_templates",
	EmailTemplates:     model.DefaultEmailTemplates,
}

// DatabaseComposer inits database stack.
type DatabaseComposer interface {
	Compose() (
		model.AppStorage,
		model.UserStorage,
		model.TokenStorage,
		model.TokenService,
		model.Encryptor,
		error,
	)
}

// NewServer creates backend service.
func NewServer(setting model.ServerSettings, db DatabaseComposer, options ...func(*Server) error) (model.Server, error) {

	appStorage, userStorage, tokenStorage, tokenService, encryptor, err := db.Compose()
	if err != nil {
		return nil, err
	}
	s := Server{AppStrg: appStorage, UserStrg: userStorage}

	sessionStorage, err := sessionStorage(setting.SessionStorage)
	if err != nil {
		return nil, err
	}
	sessionService := model.NewSessionManager(setting.SessionDuration, sessionStorage)

	ms, err := mailService(setting.MailService, setting.EmailTemplates, setting.EmailTemplatesPath)
	if err != nil {
		return nil, err
	}

	//env variable could rewrite this option
	hostName := os.Getenv("HOST_NAME")
	if len(hostName) == 0 {
		hostName = setting.Host
	}

	staticFiles := html.StaticFilesPath{
		StylesPath:  path.Join(setting.StaticFolderPath, "css"),
		ScriptsPath: path.Join(setting.StaticFolderPath, "js"),
		PagesPath:   setting.StaticFolderPath,
		ImagesPath:  path.Join(setting.StaticFolderPath, "img"),
	}

	routerSettings := web.RouterSetting{
		AppStorage:     appStorage,
		UserStorage:    userStorage,
		TokenStorage:   tokenStorage,
		TokenService:   tokenService,
		SessionService: sessionService,
		SessionStorage: sessionStorage,
		Encryptor:      encryptor,
		EmailService:   ms,
		WebRouterSettings: []func(*html.Router) error{
			html.StaticPathOptions(staticFiles),
			html.HostOption(hostName),
		},
		APIRouterSettings: []func(*api.Router) error{
			api.HostOption(hostName),
		},
		AdminRouterSettings: []func(*admin.Router) error{
			admin.HostOption(hostName),
			admin.ConfigPathOption(setting.ConfigPath),
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

//Server is DynamoDB backed server
type Server struct {
	MainRouter *web.Router
	AppStrg    model.AppStorage
	UserStrg   model.UserStorage
}

//Router return default router
func (s *Server) Router() model.Router {
	return s.MainRouter
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.MainRouter.ServeHTTP(w, r)
}

//AppStorage app storage
func (s *Server) AppStorage() model.AppStorage {
	return s.AppStrg
}

//UserStorage return user storage
func (s *Server) UserStorage() model.UserStorage {
	return s.UserStrg
}

func mailService(serviceType model.MailServiceType, templates model.EmailTemplates, templatesPath string) (model.EmailService, error) {
	tpltr, err := model.NewEmailTemplater(templates, templatesPath)
	if err != nil {
		return nil, err
	}
	switch serviceType {
	case model.MailServiceMailgun:
		return mailgun.NewEmailServiceFromEnv(tpltr)
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

//ImportApps imports apps from file
func (s *Server) ImportApps(filename string) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return s.AppStorage().ImportJSON(data)
}

//ImportUsers import users from file
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
