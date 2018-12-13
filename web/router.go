package web

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/api"
	"github.com/madappgang/identifo/web/html"
)

//RouterSetting settings for root http router
type RouterSetting struct {
	AppStorage        model.AppStorage
	UserStorage       model.UserStorage
	TokenStorage      model.TokenStorage
	TokenService      model.TokenService
	EmailService      model.EmailService
	Logger            *log.Logger
	APIRouterSettings []func(*api.Router) error
	WebRouterSettings []func(*html.Router) error
}

//NewRouter create and init root http router
func NewRouter(settings RouterSetting) (model.Router, error) {

	r := Router{}
	var err error
	r.APIRouter, err = api.NewRouter(
		settings.Logger,
		settings.AppStorage,
		settings.UserStorage,
		settings.TokenStorage,
		settings.TokenService,
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
		settings.EmailService,
		settings.WebRouterSettings...,
	)

	if err != nil {
		return nil, err
	}

	//TODO: Admin panel router
	r.APIRouterPath = "/api"
	r.WebRouterPath = "/web"
	r.AdminRouterPath = "/admin"
	r.setupRoutes()
	return &r, nil
}

//Router - root router to handle REST API, web and admin
type Router struct {
	APIRouter   model.Router
	WebRouter   model.Router
	AdminRouter model.Router
	RootRouter  *http.ServeMux

	APIRouterPath   string
	WebRouterPath   string
	AdminRouterPath string
}

//ServeHTTP identifo.Router protocol implementation
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.RootRouter.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter = http.NewServeMux()
	ar.RootRouter.Handle(ar.WebRouterPath+"/", http.StripPrefix(ar.WebRouterPath, ar.WebRouter))
	ar.RootRouter.Handle("/", ar.APIRouter)

	//TODO: add admin panel router
}
