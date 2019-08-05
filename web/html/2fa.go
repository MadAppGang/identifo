package html

import (
	"html/template"
	"net/http"
	"path"
	"time"

	ijwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"github.com/xlzd/gotp"
)

// DisableTFA handles TFA disablement form submission (POST request).
func (ar *Router) DisableTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(model.TokenContextKey).(ijwt.Token)
		if !ok {
			ar.Logger.Println("Error getting token from context")
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

// ResetTFA handles TFA resetting form submission (POST request).
func (ar *Router) ResetTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(model.TokenContextKey).(ijwt.Token)
		if !ok {
			ar.Logger.Println("Error getting token from context")
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

		totp := gotp.NewDefaultTOTP(user.TFAInfo().Secret)
		if verified := totp.Verify(r.FormValue("tfa_code"), int(time.Now().Unix())); !verified {
			SetFlash(w, FlashErrorMessageKey, "Invalid TFA code")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		if _, err := ar.UserStorage.UpdateUser(token.UserID(), user); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		successPath := path.Join(ar.PathPrefix, "tfa/reset/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}
}

// ResetTFAHandler handles reset TFA GET request.
func (ar *Router) ResetTFAHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(ar.StaticFilesPath.PagesPath, ar.StaticPages.ResetTFA))
	if err != nil {
		ar.Logger.Fatalln("Cannot parse ResetTFA template.", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		token, ok := r.Context().Value(model.TokenContextKey).(ijwt.Token)
		if !ok {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		user, err := ar.UserStorage.UserByID(token.UserID())
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		tfa := model.TFAInfo{
			IsEnabled: true,
			Secret:    gotp.RandomSecret(16),
		}
		user.SetTFAInfo(tfa)

		if _, err := ar.UserStorage.UpdateUser(user.ID(), user); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		data := map[string]interface{}{
			"Error":     errorMessage,
			"Token":     token,
			"Prefix":    ar.PathPrefix,
			"TFASecret": tfa.Secret,
		}

		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
