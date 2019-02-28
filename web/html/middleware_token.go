package html

import (
	"context"
	"net/http"
	"path"

	"github.com/madappgang/identifo/web/shared"

	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

//ResetTokenMiddleware checks token in questy and validate it
func (ar *Router) ResetTokenMiddleware() negroni.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/reset/error")
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tstr := ""
		switch r.Method {
		case http.MethodGet:
			tstr = r.URL.Query().Get("token")
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				break
			}

			tstr = r.FormValue("token")
		}

		if tstr == "" {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		_, err := shared.ParseToken(tstr, ar.TokenService, model.RefrestTokenType)
		if err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
