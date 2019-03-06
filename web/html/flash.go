package html

import (
	"net/http"
)

const (
	// FlashErrorMessageKey flash message key to keep error message across pages.
	FlashErrorMessageKey = "error"
)

// SetFlash sets new flash message
func SetFlash(w http.ResponseWriter, name, value string) {
	setCookie(w, name, value, 600)
}

// GetFlash gets flash message
func GetFlash(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	value, err := getCookie(r, name)
	deleteCookie(w, name)
	return value, err
}
