package html

import (
	"context"
	"net/http"
	"path"

	jwtService "github.com/madappgang/identifo/jwt/service"
	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

// ResetTokenMiddleware extracts reset token and validates it.
func (ar *Router) ResetTokenMiddleware() negroni.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")
	tokenValidator := jwtValidator.NewValidator(
		[]string{"identifo"},
		[]string{ar.TokenService.Issuer()},
		[]string{},
		[]string{jwtService.ResetTokenType},
	)

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tstr := ""
		switch r.Method {
		case http.MethodGet:
			tstr = r.URL.Query().Get("token")
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				break
			}
			tstr = r.FormValue("token")
		}

		if tstr == "" {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		token, err := ar.TokenService.Parse(tstr)
		if err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		if err = tokenValidator.Validate(token); err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
		ctx = context.WithValue(ctx, model.TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
