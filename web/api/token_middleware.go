package api

import (
	"context"
	"net/http"

	"github.com/madappgang/identifo/jwt"
	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
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

// Token middleware extracts token and validates it
func (ar *Router) Token(tokenType string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(rw, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App id is not in request header params.", "Token.AppFromContext")
			return
		}

		tstr := jwt.ExtractTokenFromBearerHeader(r.Header.Get(TokenHeaderKey))
		if tstr == nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token is empty or invalid.", "Token.ExtractTokenFromBearerHeader")
			return
		}

		v := jwtValidator.NewValidator(app.ID(), ar.tokenService.Issuer(), "", tokenType)
		token, err := ar.tokenService.Parse(string(tstr))
		if err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.tokenService_Parse")
			return
		}
		if err := v.Validate(token); err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, err.Error(), "Token.Validate(token)")
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
		ctx = context.WithValue(ctx, model.TokenRawContextKey, tstr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

// tokenFromContext returns token from request context.
func tokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(model.TokenContextKey).(jwt.Token)
}
