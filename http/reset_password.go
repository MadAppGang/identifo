package http

import "net/http"

func (ar *apiRouter) ResetPassword() http.HandlerFunc {
	type res struct {
		Message interface{} `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

	}
}
