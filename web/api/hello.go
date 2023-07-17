package api

import (
	"net/http"
	"time"
)

type helloResponse struct {
	Answer string    `json:"answer,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}

// HandleHello returns hello message.
func (ar *Router) Hello(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")

	hello := helloResponse{
		Answer: "Hello, my name is Aooth server",
		Date:   time.Now(),
	}
	ar.ServeJSON(w, locale, http.StatusOK, hello)
}

// Ping returns pong message.
func (ar *Router) Ping(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")

	ar.Logger.Println("trace pong handler")
	pong := helloResponse{
		Answer: "Pong!",
		Date:   time.Now(),
	}
	ar.ServeJSON(w, locale, http.StatusOK, pong)
}
