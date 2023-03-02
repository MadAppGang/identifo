package management

import (
	"net/http"
	"time"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/sig"
)

var KeyIDHeaderKey = http.CanonicalHeaderKey("X-Nl-Key-Id")

func (ar *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		if len(r.Header[KeyIDHeaderKey]) != 1 || len(r.Header[KeyIDHeaderKey][0]) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorMANoKeyID)
			return
		}

		keyID := r.Header[KeyIDHeaderKey][0]
		key, err := ar.stor.GetKey(r.Context(), keyID)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorMAErrorGettingKeyWithID, keyID, err)
			return
		}

		if key.Active == false {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorMAErrorInactiveKey)
			return
		}

		if key.ValidTill != nil && time.Now().After(*key.ValidTill) {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorMAErrorExpiredKey)
			return
		}

		err = sig.VerifySignature(r, []byte(key.Secret))
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorMAErrorInvalidSignature, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
