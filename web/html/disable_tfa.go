package html

import (
	"html/template"
	"net/http"
	"path"

	"github.com/madappgang/identifo/model"
)

// DisableTFA handles TFA disablement form submission (POST request).
func (ar *Router) DisableTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Context().Value(model.TokenRawContextKey).(string)
		token, err := ar.TokenService.Parse(tokenString)
		if err != nil {
			ar.Logger.Println("Error parsing token. ", err)
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		user, err := ar.UserStorage.UserByID(token.UserID())
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		tfa := model.TFAInfo{
			IsEnabled: false,
			Secret:    "",
		}
		user.SetTFAInfo(tfa)

		if _, err := ar.UserStorage.UpdateUser(token.UserID(), user); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		successPath := path.Join(ar.PathPrefix, "tfa/disable/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}
}

// DisableTFAHandler handles disable TFA GET request.
func (ar *Router) DisableTFAHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(ar.StaticFilesPath.PagesPath, ar.StaticPages.DisableTFA))
	if err != nil {
		ar.Logger.Fatalln("Cannot parse DisableTFA template.", err)
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
