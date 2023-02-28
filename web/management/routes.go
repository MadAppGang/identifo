package management

import (
	"log"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/model"
	"github.com/urfave/negroni"
)

// Router is a router that handles management requests.
type Router struct {
	server       model.Server
	middleware   *negroni.Negroni
	logger       *log.Logger
	router       *mux.Router
	PathPrefix   string
	Host         *url.URL
	forceRestart chan<- bool
}

type RouterSettings struct {
	Server  model.Server
	Logger  *log.Logger
	Host    *url.URL
	Prefix  string
	Restart chan<- bool
}
