package spa

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/urfave/negroni"
)

func NewRouter(setting SPASettings, logger *log.Logger) (model.Router, error) {
	ar := Router{
		Middleware: negroni.New(middleware.NewNegroniLogger(setting.Name), negroni.NewRecovery()),
		FS:         setting.FileSystem,
	}

	// Setup logger to stdout.
	if logger == nil {
		ar.Logger = log.New(os.Stdout, fmt.Sprintf("[ %s ]: ", setting.Name), log.Ldate|log.Ltime|log.Lshortfile)
	}

	ar.Middleware.UseHandler(NewSPAHandlerFunc(setting))
	return &ar, nil
}

// login app router
type Router struct {
	Logger     *log.Logger
	Middleware *negroni.Negroni
	FS         http.FileSystem
}

// ServeHTTP implements identifo.Router interface.
func (ar *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Reroute to our internal implementation.
	ar.Middleware.ServeHTTP(w, r)
}
