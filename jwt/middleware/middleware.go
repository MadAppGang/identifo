package middleware

import (
	"context"
	"net/http"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
)

const (
	// AuthorizationHeaderKey is a header name for Bearer token.
	AuthorizationHeaderKey = "Authorization"
	// TokenTypeAccess is an access token type.
	TokenTypeAccess = "access"
	// TokenTypeRefresh is a refresh token type.
	TokenTypeRefresh = "refresh"
	// AccessTokenContextKey context key to store and retreive access token
	AccessTokenContextKey = "identifo.token.access"
	// RefreshTokenContextKey context key to store and retreive refresh token
	RefreshTokenContextKey = "identifo.token.refresh"
)

//Handler is a full copy of negroni.HandlerFunc
//this is the same http.HandlerFunc, it just has one additional parameter 'next'
//next is a reference to the next handler in the handler chain
type Handler func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

//ErrorHandler interface for handling error from middleware
//rw - readwriter to write error to (JSON, HTML or other error to the client)
//errorType - error returned from middleware, you can get description by calling errorType.Description()
//status - http status code
type ErrorHandler interface {
	Error(rw http.ResponseWriter, errorType Error, status int, description string)
}

//Middleware returns middleware function you can use
func Middleware(eh ErrorHandler, c validator.Config) Handler {

	v := validator.NewValidatorWithConfig(c)
	// Middleware middleware functions extracts token and validates it and store the parsed token in the context
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tokenBytes := jwt.ExtractTokenFromBearerHeader(r.Header.Get(AuthorizationHeaderKey))
		if tokenBytes == nil {
			eh.Error(rw, ErrorTokenIsEmpty, http.StatusBadRequest, "")
			return
		}
		tokenString := string(tokenBytes)

		token, err := v.ValidateString(tokenString)
		if err != nil {
			eh.Error(rw, ErrorTokenIsInvalid, http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), model.TokenContextKey, token)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

// TokenFromContext returns token from request context.
func TokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(model.TokenContextKey).(jwt.Token)
}
