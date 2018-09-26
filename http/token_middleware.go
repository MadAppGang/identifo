package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/jwt"
	"github.com/urfave/negroni"
)

const (
	//TokenHeaderKey header key to keep Bearer token
	TokenHeaderKey = "Authorization"
	//TokenHeaderKeyPrefix token prefix regarding RFCXXX
	TokenHeaderKeyPrefix = "BEARER "
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

		tstr := extractToken(r.Header.Get(TokenHeaderKey))
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

		ctx := context.WithValue(r.Context(), TokenContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

func extractToken(token string) []byte {
	token = strings.TrimSpace(token)
	if (len(token) <= len(TokenHeaderKeyPrefix)) ||
		(strings.ToUpper(token[0:len(TokenHeaderKeyPrefix)]) != TokenHeaderKeyPrefix) {
		return nil
	}

	token = token[len(TokenHeaderKeyPrefix):]
	return []byte(token)
}
