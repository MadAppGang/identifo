package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/jwt"
	jwtValidator "github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

const (
	// TokenHeaderKey is a header name for Bearer token.
	TokenHeaderKey = "Authorization"
)

// Token middleware extracts token and validates it.
func (ar *Router) Token(tokenType model.TokenType, scopes []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			locale := r.Header.Get("Accept-Language")

			app := middleware.AppFromContext(r.Context())
			if len(app.ID) == 0 {
				ar.LocalizedError(rw, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
				return
			}

			tokenBytes := jwt.ExtractTokenFromBearerHeader(r.Header.Get(TokenHeaderKey))
			if tokenBytes == nil {
				ar.LocalizedError(rw, locale, http.StatusBadRequest, l.ErrorAPIRequestTokenInvalid)
				return
			}
			tokenString := string(tokenBytes)

			v := jwtValidator.NewValidator(
				[]string{app.ID},
				[]string{ar.server.Services().Token.Issuer()},
				nil,
				[]string{string(tokenType)},
			)
			token, err := ar.server.Services().Token.Parse(tokenString)
			if err != nil {
				ar.LocalizedError(rw, locale, http.StatusBadRequest, l.ErrorAPITokenParseError, err)
				return
			}
			if err := v.Validate(token); err != nil {
				ar.LocalizedError(rw, locale, http.StatusBadRequest, l.ErrorTokenInvalidError, err)
				return
			}

			if blacklisted := ar.server.Storages().Blocklist.IsBlacklisted(tokenString); blacklisted {
				ar.LocalizedError(rw, locale, http.StatusBadRequest, l.ErrorTokenBlocked)
				return
			}

			ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
			ctx = context.WithValue(ctx, model.TokenRawContextKey, tokenBytes)
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}
}

// tokenFromContext returns token from request context.
func tokenFromContext(ctx context.Context) *model.JWToken {
	return ctx.Value(model.TokenContextKey).(*model.JWToken)
}
