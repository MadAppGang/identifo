package html

import (
	"net/http"
	"path"

	"github.com/madappgang/identifo/model"
)

//ResetPassword handles password reset form submition
func (ar *Router) ResetPassword() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")
		if err := model.StrongPswd(password); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		tokenString := r.Context().Value(model.TokenRawContextKey).(string)
		token, _ := ar.TokenService.Parse(tokenString)

		err := ar.UserStorage.ResetPassword(token.UserID(), password)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		successPath := path.Join(".", ar.PathPrefix, "/reset/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}

}
