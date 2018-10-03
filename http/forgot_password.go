package http

import (
	"html/template"
	"net/http"
)

func (ar *apiRouter) ForgotPassword() http.HandlerFunc {
	tmpl, err := template.ParseFiles("../../tmpl/forget-password.html")

	type userRequest struct {
		Username string `json:"username,omitempty"`
	}

	type ForgotPasswordResponse struct {
		Result string `json:"result,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}

			tmpl.Execute(w, nil)
		case "POST":
			d := userRequest{}
			if ar.MustParseJSON(w, r, &d) != nil {
				return
			}

			_, err := ar.userStorage.UserByName(d.Username)
			if err != nil {
				ar.Error(w, err, http.StatusBadRequest, "")
				return
			}

			// TODO: send email with a reset password link

			result := ForgotPasswordResponse{Result: "ok"}

			ar.ServeJSON(w, http.StatusOK, result)
		}
	}
}
