package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/madappgang/identifo/mailgun"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/ses"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/html"
)

//DefaultSettings default serve settings
var DefaultSettings = model.ServerSettings{
	StaticFolderPath: "./static",
	PEMFolderPath:    "./pem",
	PrivateKey:       "private.pem",
	PublicKey:        "public.pem",
	Algorithm:        model.TokenServiceAlgorithmAuto,
	Issuer:           "identifo",
	MailService:      model.MailServiceMailgun,
}

//DatabaseComposer init database stack
type DatabaseComposer interface {
	Compose() (
		model.AppStorage,
		model.UserStorage,
		model.TokenStorage,
		model.TokenService,
		error,
	)
}

//NewServer create backend service
func NewServer(setting model.ServerSettings, db DatabaseComposer, options ...func(*Server) error) (model.Server, error) {
	s := Server{}

	appStorage, userStorage, tokenStorage, tokenService, err := db.Compose()
	if err != nil {
		return nil, err
	}
	s.AppStrg = appStorage
	s.UserStrg = userStorage

	ms, err := mailService(setting.MailService)
	if err != nil {
		return nil, err
	}

	staticFiles := html.StaticFilesPath{
		StylesPath:  path.Join(setting.StaticFolderPath, "css"),
		ScriptsPath: path.Join(setting.StaticFolderPath, "js"),
		PagesPath:   setting.StaticFolderPath,
		ImagesPath:  path.Join(setting.StaticFolderPath, "img"),
	}

	routerSettings := web.RouterSetting{
		AppStorage:   appStorage,
		UserStorage:  userStorage,
		TokenStorage: tokenStorage,
		TokenService: tokenService,
		EmailService: ms,
		WebRouterSettings: []func(*html.Router) error{
			html.StaticPathOptions(staticFiles),
		},
	}
	r, err := web.NewRouter(routerSettings)
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

func mailService(serviceType model.MailServiceType) (model.EmailService, error) {
	switch serviceType {
	case model.MailServiceMailgun:
		return mailgun.NewEmailServiceFromEnv(), nil
	case model.MailServiceAWS:
		return ses.NewEmailServiceFromEnv()
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
