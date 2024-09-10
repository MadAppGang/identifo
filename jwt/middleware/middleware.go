package middleware

import (
	"context"
	"net/http"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/jwt/validator"
	"github.com/madappgang/identifo/v2/model"
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

// Handler is a full copy of negroni.HandlerFunc
// this is the same http.HandlerFunc, it just has one additional parameter 'next'
// next is a reference to the next handler in the handler chain
type Handler func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// JWTErrorHandlerContext is a context for handling errors from JWTv2 middleware.
type JWTErrorHandlerContext struct {
	// ResponseWriter is an http.ResponseWriter to write error to (JSON, HTML or other error to the client).
	ResponseWriter http.ResponseWriter
	// Request is an http.Request to get more information about the request.
	Request *http.Request
	// ErrorType is an error returned from middleware, you can get description by calling errorType.Description().
	ErrorType Error
	// Status is an http status code.
	Status int
	// Description is a description of the error.
	Description string
	// Token is a parsed but not valid token, may be nil.
	Token model.Token
}

// JWTErrorHandler is an interface for handling errors from JWTv2 middleware.
type JWTErrorHandler interface {
	Error(errContext *JWTErrorHandlerContext)
}

// JWTv2 returns middleware function you can use to handle JWT token auth.
// It extracts token from Authorization header, validates it and stores the parsed token in the context,
// use TokenFromContext to get token from context.
// It uses JWTErrorHandler to handle errors.
func JWTv2(eh JWTErrorHandler, c validator.Config) (Handler, error) {
	v, err := validator.NewValidatorWithConfig(c)
	if err != nil {
		return nil, err
	}

	reportError := func(
		rw http.ResponseWriter,
		r *http.Request,
		errType Error,
		status int,
		description string,
		token model.Token,
	) {
		ec := &JWTErrorHandlerContext{
			ResponseWriter: rw,
			Request:        r,
			ErrorType:      errType,
			Status:         status,
			Description:    description,
			Token:          token,
		}

		eh.Error(ec)
	}

	// Middleware middleware functions extracts token and validates it and store the parsed token in the context
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tokenBytes := jwt.ExtractTokenFromBearerHeader(r.Header.Get(AuthorizationHeaderKey))
		if tokenBytes == nil {
			reportError(rw, r, ErrorTokenIsEmpty, http.StatusBadRequest, "", nil)
			return
		}
		tokenString := string(tokenBytes)

		token, err := v.ValidateString(tokenString)
		if err != nil {
			invalidToken, _ := jwt.ParseTokenString(tokenString)
			reportError(rw, r, ErrorTokenIsInvalid, http.StatusBadRequest, err.Error(), invalidToken)
			return
		}

		r = appendRequestContextValue(r, model.TokenContextKey, token)

		next.ServeHTTP(rw, r)
	}, nil
}

func appendRequestContextValue(r *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

// ErrorHandler interface for handling error from middleware
// rw - http.ResponseWriter to write error to (JSON, HTML or other error to the client)
// errorType - error returned from middleware, you can get description by calling errorType.Description()
// status - http status code
type ErrorHandler interface {
	Error(rw http.ResponseWriter, errorType Error, status int, description string)
}

// jwtErrorHandler is a wrapper for ErrorHandler to make it compatible with JWTErrorHandler
type jwtErrorHandler struct {
	eh ErrorHandler
}

func (j jwtErrorHandler) Error(errContext *JWTErrorHandlerContext) {
	j.eh.Error(errContext.ResponseWriter, errContext.ErrorType, errContext.Status, errContext.Description)
}

// JWT returns middleware function you can use to handle JWT token auth
// Deprecated: use JWT instead
func JWT(eh ErrorHandler, c validator.Config) (Handler, error) {
	return JWTv2(jwtErrorHandler{eh}, c)
}

// TokenFromContext returns token from request context.
// Or nil if there is no token in context.
func TokenFromContext(ctx context.Context) model.Token {
	v, _ := ctx.Value(model.TokenContextKey).(model.Token)
	return v
}
