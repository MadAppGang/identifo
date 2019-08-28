package html

import (
	"net/http"
	"path"
	"time"

	ijwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

// DisableTFA handles TFA disablement form submission (POST request).
func (ar *Router) DisableTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Logger.Println("Error getting token from context")
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}
		tokenString := string(tokenBytes)

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

		// Invalidate reset token after use.
		if err := ar.TokenBlacklist.Add(tokenString); err != nil {
			ar.Logger.Printf("Cannot blacklist reset token after use: %s\n", err)
		}

		successPath := path.Join(ar.PathPrefix, "tfa/disable/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}
}

// DisableTFAHandler handles disable TFA GET request.
func (ar *Router) DisableTFAHandler() http.HandlerFunc {
	//tmpl, err := template.ParseFiles(path.Join(ar.StaticFilesPath.PagesPath, ar.StaticPages.DisableTFA))
	tmpl, err := ar.StaticFilesStorage.ParseTemplate(model.DisableTFATemplateName)
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
		tokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Logger.Println("Error getting token bytes from context")
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}
		tokenString := string(tokenBytes)

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

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Logger.Println("Error getting app from context")
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		tfaCode := r.FormValue("tfa_code")
		if len(tfaCode) == 0 {
			SetFlash(w, FlashErrorMessageKey, "Empty TFA code")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		totp := gotp.NewDefaultTOTP(user.TFAInfo().Secret)
		dontNeedVerification := app.DebugTFACode() != "" && tfaCode == app.DebugTFACode()

		if verified := totp.Verify(tfaCode, int(time.Now().Unix())); !(verified || dontNeedVerification) {
			SetFlash(w, FlashErrorMessageKey, "Invalid TFA code")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		if _, err := ar.UserStorage.UpdateUser(token.UserID(), user); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		// Invalidate reset token after use.
		if err := ar.TokenBlacklist.Add(tokenString); err != nil {
			ar.Logger.Printf("Cannot blacklist reset token after use: %s\n", err)
		}

		successPath := path.Join(ar.PathPrefix, "tfa/reset/success")
		http.Redirect(w, r, successPath, http.StatusMovedPermanently)
	}
}

// ResetTFAHandler handles reset TFA GET request.
func (ar *Router) ResetTFAHandler() http.HandlerFunc {
	//tmpl, err := template.ParseFiles(path.Join(ar.StaticFilesPath.PagesPath, ar.StaticPages.ResetTFA))
	tmpl, err := ar.StaticFilesStorage.ParseTemplate(model.ResetTFATemplateName)
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
