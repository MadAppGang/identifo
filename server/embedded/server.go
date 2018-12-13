package embedded

import (
	"net/http"
	"path"

	"github.com/madappgang/identifo/boltdb"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mailgun"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/html"
)

//Settings is extended settings for DynamoDB serve
type Settings struct {
	model.ServerSettings
	DBPath string
}

//DefaultSettings default serve settings
var DefaultSettings = Settings{
	ServerSettings: model.ServerSettings{
		StaticFolderPath: "./static",
		PEMFolderPath:    "./pem",
		PrivateKey:       "private.pem",
		PublicKey:        "public.pem",
		Algorithm:        model.TokenServiceAlgorithmAuto,
		Issuer:           "identifo",
		MailService:      model.MailServiceMailgun,
	},
	DBPath: "db.db",
}

//NewServer create DynamoDB backend service
func NewServer(setting Settings, options ...func(*Server) error) (model.Server, error) {
	s := Server{}

	db, err := boltdb.InitDB(setting.DBPath)
	if err != nil {
		return nil, err
	}
	appStorage, err := boltdb.NewAppStorage(db)
	if err != nil {
		return nil, err
	}
	s.AppStrg = appStorage

	userStorage, err := boltdb.NewUserStorage(db)
	if err != nil {
		return nil, err
	}
	s.UserStrg = userStorage

	tokenStorage, err := boltdb.NewTokenStorage(db)
	if err != nil {
		return nil, err
	}

	tokenService, err := jwt.NewTokenService(
		path.Join(setting.PEMFolderPath, setting.PrivateKey),
		path.Join(setting.PEMFolderPath, setting.PublicKey),
		setting.Issuer,
		setting.Algorithm,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, err
	}

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
	default:
		return nil, model.ErrorNotImplemented
	}
}
