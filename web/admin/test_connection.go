package admin

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
)

// TestConnection validates different connection types, if server could connect to that.
func (ar *Router) TestConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		tc := model.TestConnection{}
		if err := ar.mustParseJSON(w, r, &tc); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIJsonParseError, err.Error())
			return
		}

		tester, err := storage.NewConnectionTester(tc)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorStorageVerificationFindError, err.Error())
			return
		}
		if tester == nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorStorageUnsupportedType)
			return
		}
		err = tester.Connect()
		if err != nil {
			if err != nil {
				ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorStorageVerificationFindError, err.Error())
				return
			}
		}
		ar.ServeJSON(w, locale, http.StatusOK, tc)
	}
}
