package api

import (
	"net/http"
	"time"
)

//HandleHello - returns hello message
func (ar *Router) HandleHello() http.HandlerFunc {
	type helloResponse struct {
		Answer string    `json:"answer,omitempty"`
		Date   time.Time `json:"date,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ar.logger.Println("trace Hello handler")
		hello := helloResponse{
			Answer: "Hello, my name is Identifo",
			Date:   time.Now(),
		}
		ar.ServeJSON(w, http.StatusOK, hello)
	}
}

// HandlePing returns pong message.
func (ar *Router) HandlePing() http.HandlerFunc {
	type pongResponse struct {
		Message string    `json:"message,omitempty"`
		Date    time.Time `json:"date,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ar.logger.Println("trace pong handler")
		pong := pongResponse{
			Message: "Pong!",
			Date:    time.Now(),
		}
		ar.ServeJSON(w, http.StatusOK, pong)
	}
}
