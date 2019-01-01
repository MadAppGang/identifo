package html

import (
	"context"
	"net/http"
	"path"

	"github.com/madappgang/identifo/jwt"
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

		v := jwt.NewValidator("identifo", ar.TokenService.Issuer(), "")
		token, err := ar.TokenService.Parse(string(tstr))
		if err != nil {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		if err := v.Validate(token); err != nil {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		if model.ResetTokenType != token.Type() {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
