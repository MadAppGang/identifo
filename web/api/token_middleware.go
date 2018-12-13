package api

import (
	"context"
	"net/http"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

const (
	//TokenHeaderKey header key to keep Bearer token
	TokenHeaderKey = "Authorization"
	//TokenTypeRefresh is to handle refresh as bearer token
	TokenTypeRefresh = "refresh"
	//TokenTypeAccess is to handle access token type as bearer token
	TokenTypeAccess = "access"
)

//Token middleware extracts token and validates it
func (ar *Router) Token(tokenType string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		app := appFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(rw, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		tstr := jwt.ExtractTokenFromBearerHeader(r.Header.Get(TokenHeaderKey))
		if tstr == nil {
			ar.Error(rw, ErrorRequestInvalidToken, http.StatusBadRequest, "")
			return
		}
		v := jwt.NewValidator(app.ID(), ar.tokenService.Issuer(), "")
		token, err := ar.tokenService.Parse(string(tstr))
		if err != nil {
			ar.Error(rw, ErrorRequestInvalidToken, http.StatusBadRequest, "")
			return
		}
		if err := v.Validate(token); err != nil {
			ar.Error(rw, ErrorRequestInvalidToken, http.StatusBadRequest, err.Error())
			return
		}
		if tokenType != token.Type() {
			ar.Error(rw, ErrorRequestInvalidToken, http.StatusBadRequest, "Invalid token type")
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
		ctx = context.WithValue(r.Context(), model.TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}
