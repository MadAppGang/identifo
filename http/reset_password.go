package http

import (
	"net/http"
)

func (ar *apiRouter) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")
		if err := StrongPswd(password); err != nil {
			SetFlash(w, ErrorMessageKey, err.Error())
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		tokenString := r.Context().Value(TokenRawContextKey).(string)
		token, _ := ar.tokenService.Parse(tokenString)

		err := ar.userStorage.ResetPassword(token.UserID(), password)
		if err != nil {
			SetFlash(w, ErrorMessageKey, "Server Error")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		http.Redirect(w, r, "./reset/success", http.StatusMovedPermanently)
	}
}
