package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/jwt"
	jwtValidator "github.com/madappgang/identifo/v2/jwt/validator"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

const (
	// TokenHeaderKey is a header name for Bearer token.
	TokenHeaderKey = "Authorization"
)

// Token middleware extracts token and validates it.
func (ar *Router) Token(tokenType string, scopes []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			locale := r.Header.Get("Accept-Language")

			app := middleware.AppFromContext(r.Context())
			if len(app.ID) == 0 {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
				return
			}

			tokenBytes := jwt.ExtractTokenFromBearerHeader(r.Header.Get(TokenHeaderKey))
			if tokenBytes == nil {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIRequestTokenInvalid)
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
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPITokenParseError, err)
				return
			}
			if err := v.Validate(token); err != nil {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorTokenInvalidError, err)
				return
			}

			if blacklisted := ar.server.Storages().Blocklist.IsBlacklisted(tokenString); blacklisted {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorTokenBlocked)
				return
			}

			if len(scopes) > 0 {
				ts := strings.Split(token.Scopes(), " ")
				if len(model.SliceIntersect(ts, scopes)) == 0 {
					ar.Error(rw, locale, http.StatusUnauthorized, l.ErrorAPPLoginNoScope)
					return
				}
			}

			ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
			ctx = context.WithValue(ctx, model.TokenRawContextKey, tokenBytes)
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}
}

// tokenFromContext returns token from request context.
func tokenFromContext(ctx context.Context) model.Token {
	return ctx.Value(model.TokenContextKey).(model.Token)
}
