package web

import (
	"log"
	"net/http"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/admin"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
)

// RouterSetting contains settings for root http router.
type RouterSetting struct {
	AppStorage              model.AppStorage
	UserStorage             model.UserStorage
	TokenStorage            model.TokenStorage
	VerificationCodeStorage model.VerificationCodeStorage
	TokenService            jwtService.TokenService
	SMSService              model.SMSService
	EmailService            model.EmailService
	SessionService          model.SessionService
	SessionStorage          model.SessionStorage
	ConfigurationStorage    model.ConfigurationStorage
	Logger                  *log.Logger
	APIRouterSettings       []func(*api.Router) error
	WebRouterSettings       []func(*html.Router) error
	AdminRouterSettings     []func(*admin.Router) error
}

// NewRouter creates and inits root http router.
func NewRouter(settings RouterSetting) (model.Router, error) {
	r := Router{}
	var err error

	r.APIRouter, err = api.NewRouter(
		settings.Logger,
		settings.AppStorage,
		settings.UserStorage,
		settings.TokenStorage,
		settings.VerificationCodeStorage,
		settings.TokenService,
		settings.SMSService,
		settings.EmailService,
		settings.APIRouterSettings...,
	)
	if err != nil {
		return nil, err
	}

	r.WebRouter, err = html.NewRouter(
		settings.Logger,
		settings.AppStorage,
		settings.UserStorage,
		settings.TokenStorage,
		settings.TokenService,
		settings.SMSService,
		settings.EmailService,
		settings.WebRouterSettings...,
	)

	if err != nil {
		return nil, err
	}

	r.AdminRouter, err = admin.NewRouter(
		settings.Logger,
		settings.SessionService,
		settings.SessionStorage,
		settings.AppStorage,
		settings.UserStorage,
		settings.ConfigurationStorage,
		settings.AdminRouterSettings...,
	)

	if err != nil {
		return nil, err
	}

	r.APIRouterPath = "/api"
	r.WebRouterPath = "/web"
	r.AdminRouterPath = "/admin"

	r.setupRoutes()
	return &r, nil
}

// Router is a root router to handle REST API, web, and admin requests.
type Router struct {
	APIRouter        model.Router
	WebRouter        model.Router
	AdminRouter      model.Router
	AdminPanelRouter model.Router
	RootRouter       *http.ServeMux

	APIRouterPath        string
	WebRouterPath        string
	AdminRouterPath      string
	AdminPanelRouterPath string
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.RootRouter.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter = http.NewServeMux()
	ar.RootRouter.Handle("/", ar.APIRouter)
	ar.RootRouter.Handle(ar.WebRouterPath+"/", http.StripPrefix(ar.WebRouterPath, ar.WebRouter))
	ar.RootRouter.Handle(ar.AdminRouterPath+"/", http.StripPrefix(ar.AdminRouterPath, ar.AdminRouter))
}
