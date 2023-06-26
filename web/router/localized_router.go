package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	l "github.com/madappgang/identifo/v2/localization"
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
		ar.Error(w, locale, http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(data); err != nil {
		log.Printf("error writing http response: %s", err)
	}
}

func NewLocalizedError(status int, locale string, errID l.LocalizedString, details ...any) *LocalizedError {
	return &LocalizedError{
		status:  status,
		locale:  locale,
		errID:   errID,
		details: details,
	}
}

type LocalizedError struct {
	status  int
	locale  string
	errID   l.LocalizedString
	details []any
}

func (e *LocalizedError) Error() string {
	return fmt.Sprintf("localized error: %v (status=%v). Details: %v.", e.errID, e.status, e.details)
}

func (ar *LocalizedRouter) ErrorResponse(w http.ResponseWriter, err error) {
	ar.Logger.Printf("api error: %v", err)

	switch e := err.(type) {
	case *LocalizedError:
		ar.error(w, 3, e.locale, e.status, e.errID, e.details...)
	default:
		ar.error(w, 3, "", http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
	}
}

// Error writes an API error message to the response and logger.
func (ar *LocalizedRouter) Error(w http.ResponseWriter, locale string, status int, errID l.LocalizedString, details ...any) {
	ar.error(w, 2, locale, status, errID, details...)
}

func (ar *LocalizedRouter) error(w http.ResponseWriter, callerDepth int, locale string, status int, errID l.LocalizedString, details ...any) {
	// errorResponse is a generic response for sending an error.
	type errorResponse struct {
		ID       string `json:"id"`
		Message  string `json:"message,omitempty"`
		Status   int    `json:"status"`
		Location string `json:"location"`
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
		},
	}

	encodeErr := json.NewEncoder(w).Encode(resp)
	if encodeErr != nil {
		ar.Logger.Printf("error writing http response: %s", errID)
	}
}
