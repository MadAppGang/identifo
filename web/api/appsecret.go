package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
)

const (
	// SignatureHeaderKey header stores HMAC signature digest.
	SignatureHeaderKey = "Digest"
	// SignatureHeaderValuePrefix is a signature prefix, indicating hash algorithm, hardcoded for now, could be dynamic in the future.
	SignatureHeaderValuePrefix = "SHA-256="
	// TimestampHeaderKey header stores timestamp.
	TimestampHeaderKey = "X-Identifo-Timestamp"
)

// SignatureHandler returns middleware that handles request signature.
// More info: https://identifo.madappgang.com/#ca6498ab-b3dc-4c1e-a5b0-2dd633831e2d.
func (ar *Router) SignatureHandler() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		var body []byte
		t := r.Header.Get(TimestampHeaderKey)

		if r.Method == "GET" {
			body = []byte(r.URL.RequestURI() + t)
			ar.logger.Info("RequestURI to sign (GET request)",
				"uriData", r.URL.RequestURI()+t)
		} else {
			// Extract body.
			b, err := io.ReadAll(r.Body)
			if err != nil {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
				return
			}

			if len(b) == 0 {
				b = []byte(r.URL.RequestURI() + t)
				ar.logger.Info("RequestURI to sign (POST request)",
					"uriData", r.URL.RequestURI()+t)
			}

			body = b
		}

		if app.Type != model.Web {
			// Read request signature from header and decode it.
			reqMAC := extractSignature(r.Header.Get(SignatureHeaderKey))
			if reqMAC == nil {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIRequestSignatureInvalid)
				return

			}
			if err := validateBodySignature(body, reqMAC, []byte(app.Secret)); err != nil {
				ar.Error(rw, locale, http.StatusBadRequest, l.ErrorAPIRequestSignatureValidationError, err)
				return
			}
		}

		if r.Method != "GET" && r.Body != http.NoBody {
			// Return body as Reader to next handlers.
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		// Call next handler.
		next(rw, r)
	}
}

// extractSignature extracts signature from raw header value and returns its byte representation.
// Returns nil slice if something goes wrong.
func extractSignature(b64 string) []byte {
	b64 = strings.TrimSpace(b64)

	if (len(b64) <= len(SignatureHeaderValuePrefix)) ||
		(strings.ToUpper(b64[0:len(SignatureHeaderValuePrefix)]) != SignatureHeaderValuePrefix) {
		return nil
	}
	// Extract Base64 part of the signature, trim prefix.
	b64 = b64[len(SignatureHeaderValuePrefix):]

	// Decode to byte slice.
	reqMAC, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil
	}
	return reqMAC
}

// validateBodySignature checks if signature for the given request `body` matches the signature `reqMAC`, signed with `secret`.
func validateBodySignature(body, reqMAC, secret []byte) error {
	mac := hmac.New(sha256.New, secret)
	if _, err := mac.Write(body); err != nil {
		return err
	}

	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(reqMAC, expectedMAC) {
		return errors.New("request hmac is not equal to expected. ")
	}
	return nil
}
