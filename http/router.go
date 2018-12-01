package http

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/madappgang/identifo/http/api"
	"github.com/madappgang/identifo/model"
)

//RouterSetting settings for root http router
type RouterSetting struct {
	AppStorage   model.AppStorage
	UserStorage  model.UserStorage
	TokenStorage model.TokenStorage
	TokenService model.TokenService
	EmailService model.EmailService
	Logger       *log.Logger
}

//NewRouter create and init root http router
func NewRouter(settings RouterSetting) model.Router {

	r := Router{}
	r.APIRouter = api.NewRouter(
		settings.Logger,
		settings.AppStorage,
		settings.UserStorage,
		settings.TokenStorage,
		settings.TokenService,
		settings.EmailService,
		nil
	)

	r.WebRouter = html.NewRouter(
		settings.Logger,
		settings.AppStorage,
		settings.UserStorage,
		settings.TokenStorage,
		settings.TokenService,
		settings.EmailService,
		nil
	)

	
	//TODO: Admin panel router
	r.APIRouterPath = "/api"
	r.WebRouterPath = "/w"
	r.AdminRouterPath = "/admin"
	r.setupRoutes()
	return &r
}

//Router - root router to handle REST API, web and admin
type Router struct {
	APIRouter   model.Router
	WebRouter   model.Router
	AdminRouter model.Router
	RootRouter *http.ServeMux
	
	APIRouterPath  string
	WebRouterPath string
	AdminRouterPath string
}

//ServeHTTP identifo.Router protocol implementation
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//reroute to our internal implementation
	ar.Router.ServeHTTP(w, r)
}

func (ar *Router) setupRoutes() {
	ar.RootRouter := http.NewServeMux()
	topMux.Handle(ar.APIRouterPath+"/", http.StripPrefix(ar.APIRouterPath, ar.APIRouter))
	topMux.Handle(ar.WebRouterPath+"/", http.StripPrefix(ar.WebRouterPath, ar.WebRouter))
	//TODO: add admin panel router
}

