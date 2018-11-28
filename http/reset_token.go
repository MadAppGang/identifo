package http

import (
	"net/http"
)

func (ar *apiRouter) SendResetToken() http.HandlerFunc {
	type res struct {
		Message string `json:"message,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.FormValue("email")
		userExists := ar.userStorage.UserExists(email)
		if !userExists {
			// TODO: show error "there is no user with this email"
			ar.ServeJSON(w, 400, res{Message: "there is no user with this email"})
			return
		}

		// TODO: generate token

		// TODO: create message body

		_, _, err := ar.emailService.SendMessage("Reset Password", "Message Body", email)
		if err != nil {
			// TODO: show error "internal server error"
			ar.ServeJSON(w, 400, res{Message: "Error sending email"})
			return
		}

		// TODO: redirect to the success send email page
		ar.ServeJSON(w, 200, res{Message: "success"})
	}
}
