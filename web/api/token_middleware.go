package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/jwt"
	jwtValidator "github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
)

const (
	// TokenHeaderKey is a header name for Bearer token.
	TokenHeaderKey = "Authorization"
)

// Token middleware extracts token and validates it.
func (ar *Router) Token(tokenType string, scopes []string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.logger.Println("Error getting App")
			ar.Error(rw, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App id is not in request header params.", "Token.AppFromContext")
			return
		}

		tokenBytes := jwt.ExtractTokenFromBearerHeader(r.Header.Get(TokenHeaderKey))
		if tokenBytes == nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token is empty or invalid.", "Token.ExtractTokenFromBearerHeader")
			return
		}
		tokenString := string(tokenBytes)

		v := jwtValidator.NewValidator(
			[]string{app.ID, "identifo"},
			[]string{ar.server.Services().Token.Issuer()},
			[]string{},
			[]string{tokenType},
		)
		token, err := ar.server.Services().Token.Parse(tokenString)
		if err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.tokenService_Parse")
			return
		}
		if err := v.Validate(token); err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, err.Error(), "Token.Validate(token)")
			return
		}

		if blacklisted := ar.server.Storages().Blocklist.IsBlacklisted(tokenString); blacklisted {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.IsBlacklisted")
			return
		}

		if len(scopes) > 0 {
			ts := strings.Split(token.Scopes(), " ")
			if len(model.SliceIntersect(ts, scopes)) == 0 {
				ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.ScopeIsNotAllowed")
				return
			}
		}

		ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
		ctx = context.WithValue(ctx, model.TokenRawContextKey, tokenBytes)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

// tokenFromContext returns token from request context.
func tokenFromContext(ctx context.Context) model.Token {
	return ctx.Value(model.TokenContextKey).(model.Token)
}
