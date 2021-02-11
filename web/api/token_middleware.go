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
	// TokenTypeAccess is an access token type.
	TokenTypeAccess = "access"
	// TokenTypeRefresh is a refresh token type.
	TokenTypeRefresh = "refresh"
)

// Token middleware extracts token and validates it.
func (ar *Router) Token(tokenType string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
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
			[]string{app.ID()},
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

		if strings.Trim(r.RequestURI, "/ ") != "auth/tfa/finalize" {
			if payload := token.Payload(); payload != nil && payload["tfa_authorized"] == "false" {
				ar.Error(rw, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "", "Token.IsTFAuthorized")
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
func tokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(model.TokenContextKey).(jwt.Token)
}
