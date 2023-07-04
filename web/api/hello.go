package api

import (
	"net/http"
	"time"
)

// HandleHello returns hello message.
func (ar *Router) HandleHello() http.HandlerFunc {
	type helloResponse struct {
		Answer string    `json:"answer,omitempty"`
		Date   time.Time `json:"date,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		ar.Logger.Println("trace Hello handler")
		hello := helloResponse{
			Answer: "Hello, my name is Identifo",
			Date:   time.Now(),
		}
		ar.ServeJSON(w, locale, http.StatusOK, hello)
	}
}

// HandlePing returns pong message.
func (ar *Router) HandlePing(w http.ResponseWriter, r *http.Request) {
	type pongResponse struct {
		Message string    `json:"message,omitempty"`
		Date    time.Time `json:"date,omitempty"`
	}

	locale := r.Header.Get("Accept-Language")

	ar.Logger.Println("trace pong handler")
	pong := pongResponse{
		Message: "Pong!",
		Date:    time.Now(),
	}
	ar.ServeJSON(w, locale, http.StatusOK, pong)
}
