package spa

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

func NewRouter(setting SPASettings, logger *log.Logger) (model.Router, error) {
	ar := Router{
		Middleware: negroni.New(negroni.NewLogger(), negroni.NewRecovery()),
		Router:     mux.NewRouter(),
		FS:         setting.FileSystem,
	}

	// Setup logger to stdout.
	if logger == nil {
		ar.Logger = log.New(os.Stdout, "HTML_ROUTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// setup the only route we have
	ar.Router.HandleFunc("/", NewSPAHandlerFunc(setting))

	ar.Middleware.UseHandler(ar.Router)
	return &ar, nil
}

// login app router
type Router struct {
	Logger     *log.Logger
	Middleware *negroni.Negroni
	Router     *mux.Router
	FS         http.FileSystem
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.Middleware.ServeHTTP(w, r)
}
