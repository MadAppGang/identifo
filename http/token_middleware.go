package http

import (
	"context"
	"net/http"
	"path"

	"github.com/madappgang/identifo/model"

	"github.com/madappgang/identifo/jwt"
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
func (ar *apiRouter) Token(tokenType string) negroni.HandlerFunc {
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

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
		ctx = context.WithValue(r.Context(), TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

func (ar *apiRouter) ResetToken() negroni.HandlerFunc {
	errorPath := path.Join("password", "error")
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tstr := ""
		switch r.Method {
		case http.MethodGet:
			tstr = r.URL.Query().Get("token")
		case http.MethodPost:
			err := r.ParseForm()
			if err == nil {
				break
			}

			tstr = r.FormValue("token")
		}

		if tstr == "" {
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		v := jwt.NewValidator("identifo", ar.tokenService.Issuer(), "")
		token, err := ar.tokenService.Parse(string(tstr))
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

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
		ctx = context.WithValue(r.Context(), TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
