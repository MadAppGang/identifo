package html

import (
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
)

const redirectURLParam = "redirectUrl"

// ResetPassword handles password reset form submission (POST request).
func (ar *Router) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")
		if err := model.StrongPswd(password); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		tokenString := r.Context().Value(model.TokenRawContextKey).(string)
		token, err := ar.Server.Services().Token.Parse(tokenString)
		if err != nil {
			ar.Logger.Println("Error parsing token. ", err)
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		if err = ar.Server.Storages().User.ResetPassword(token.UserID(), password); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		redirectURL := strings.TrimSpace(r.URL.Query().Get(redirectURLParam))
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			return
		}

		successPath := path.Join(ar.PathPrefix, "password/reset/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}
}

// ResetPasswordHandler handles reset password GET request.
func (ar *Router) ResetPasswordHandler() http.HandlerFunc {
	tmpl, err := ar.Server.Storages().Static.ParseTemplate(model.StaticPagesNames.ResetPassword)
	if err != nil {
		ar.Logger.Fatalln("Cannot parse ResetPassword template.", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		token := r.Context().Value(model.TokenRawContextKey)
		data := map[string]interface{}{
			"Error":  errorMessage,
			"Token":  token,
			"Prefix": ar.PathPrefix,
		}

		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
