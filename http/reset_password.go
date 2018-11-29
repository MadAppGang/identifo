package http

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
)

func (ar *apiRouter) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")
		if err := StrongPswd(password); err != nil {
			SetFlash(w, ErrorMessageKey, err.Error())
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		}

		userID := r.Context().Value(TokenContextKey).(model.Token).UserID()
		fmt.Println(userID)
	}
}
