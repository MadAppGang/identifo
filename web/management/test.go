package management

import (
	"net/http"
	"time"
)

// HandlePing returns pong message.
func (ar *Router) test(w http.ResponseWriter, r *http.Request) {
	type pongResponse struct {
		Message string    `json:"message,omitempty"`
		Date    time.Time `json:"date,omitempty"`
	}

	locale := r.Header.Get("Accept-Language")

	ar.logger.Debug("trace pong handler")

	pong := pongResponse{
		Message: "Pong!",
		Date:    time.Now(),
	}
	ar.ServeJSON(w, locale, http.StatusOK, pong)
}
