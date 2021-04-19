package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/jwt"
	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
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
			[]string{app.ID},
			[]string{ar.tokenService.Issuer()},
			[]string{},
			[]string{tokenType},
		)
		token, err := ar.tokenService.Parse(tokenString)
		if err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.tokenService_Parse")
			return
		}
		if err := v.Validate(token); err != nil {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, err.Error(), "Token.Validate(token)")
			return
		}

		if blacklisted := ar.tokenBlacklist.IsBlacklisted(tokenString); blacklisted {
			ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.IsBlacklisted")
			return
		}

		if len(scopes) > 0 {
			ts := strings.Split(token.Scopes(), " ")
			if len(intersect(ts, scopes)) == 0 {
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

// simple intersection of two slices, with complexity: O(n^2)
// there is better algorithms around, this one is simple and scopes are usually 1-3 items in it
func intersect(a, b []string) []string {
	res := make([]string, 0)

	for _, e := range a {
		if contains(b, e) {
			res = append(res, e)
		}
	}

	return res
}

// tokenFromContext returns token from request context.
func tokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(model.TokenContextKey).(jwt.Token)
}
