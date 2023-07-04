package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/madappgang/identifo/v2/l"
)

// LocalizedRouter router with provide localized error output.
type LocalizedRouter struct {
	Logger *log.Logger
	LP     *l.Printer // localized string
}

// ServeJSON sends status code, headers and data and send it back to the user
func (ar *LocalizedRouter) ServeJSON(w http.ResponseWriter, locale string, status int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ar.LocalizedError(w, locale, http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(data); err != nil {
		log.Printf("error writing http response: %s", err)
	}
}

func (ar *LocalizedRouter) HTTPError(w http.ResponseWriter, err error, status int) {
	var le l.LocalizedError
	if errors.As(err, &le) {
		e := l.HTTPLocalizedError{
			LE:     le,
			Time:   time.Now(),
			Status: status,
		}
		ar.Error(w, e)
		return
	}

	ar.error(w, 3, "", http.StatusInternalServerError, time.Now(), l.APIInternalServerErrorWithError, err)
}

func (ar *LocalizedRouter) Error(w http.ResponseWriter, err error) {
	ar.Logger.Printf("api error: %v", err)

	var he l.HTTPLocalizedError
	if errors.As(err, &he) {
		ar.error(w, 3, he.LE.Locale, he.Status, he.Time, he.LE.ErrID, he.LE.Details...)
		return
	}

	var le l.LocalizedError
	if errors.As(err, &le) {
		ar.error(w, 3, le.Locale, 0, time.Now(), le.ErrID, le.Details...)
		return
	}

	// default unsupported error
	ar.error(w, 3, "", http.StatusInternalServerError, time.Now(), l.APIInternalServerErrorWithError, err)
}

// LocalizedError writes an API error message to the response and logger.
func (ar *LocalizedRouter) LocalizedError(w http.ResponseWriter, locale string, status int, errID l.LocalizedString, details ...any) {
	ar.error(w, 2, locale, status, time.Now(), errID, details...)
}

func (ar *LocalizedRouter) error(w http.ResponseWriter, callerDepth int, locale string, status int, t time.Time, errID l.LocalizedString, details ...any) {
	// errorResponse is a generic response for sending an error.
	type errorResponse struct {
		ID       string    `json:"id"`
		Message  string    `json:"message,omitempty"`
		Status   int       `json:"status"`
		Location string    `json:"location"`
		Time     time.Time `json:"time"`
	}

	if errID == "" {
		errID = l.APIInternalServerError
	}

	_, file, no, ok := runtime.Caller(callerDepth)
	if !ok {
		file = "unknown file"
		no = -1
	}
	message := ar.LP.SL(locale, errID, details...)

	// Log error.
	ar.Logger.Printf("api error: %v (status=%v). Details: %v. Where: %v:%d.", errID, status, message, file, no)

	// Write generic error response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := struct {
		Error errorResponse `json:"error"`
	}{
		Error: errorResponse{
			ID:       string(errID),
			Message:  message,
			Status:   status,
			Location: fmt.Sprintf("%s:%d", file, no),
			Time:     t,
		},
	}

	encodeErr := json.NewEncoder(w).Encode(resp)
	if encodeErr != nil {
		ar.Logger.Printf("error writing http response: %s", errID)
	}
}
