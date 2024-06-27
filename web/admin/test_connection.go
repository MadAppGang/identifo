package admin

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
)

// TestConnection validates different connection types, if server could connect to that.
func (ar *Router) TestConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc := model.TestConnection{}
		if err := ar.mustParseJSON(w, r, &tc); err != nil {
			ar.Error(w, fmt.Errorf("error parsing connection settings: %v", err), http.StatusBadRequest, "")
			return
		}

		tester, err := storage.NewConnectionTester(tc)
		if err != nil {
			ar.Error(w, fmt.Errorf("error creating connection tester: %v", err), http.StatusBadRequest, "")
			return
		}
		if tester == nil {
			ar.Error(w, fmt.Errorf("error creating connection tester, unsupported test type: %v", tc.Type), http.StatusBadRequest, "")
			return
		}
		err = tester.Connect()
		if err != nil {
			ar.Error(w, fmt.Errorf("connection error: %v", err), http.StatusBadRequest, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, tc)
	}
}
