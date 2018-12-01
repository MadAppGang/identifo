package html

import (
	"encoding/base64"
	"net/http"
	"time"
)

const (
	//FlashErrorMessageKey flash message key to keep error message across pages
	FlashErrorMessageKey = "error"
)

// SetFlash sets new flash message
func SetFlash(w http.ResponseWriter, name, value string) {
	c := &http.Cookie{Name: name, Value: encode(value), MaxAge: 600}
	http.SetCookie(w, c)
}

// GetFlash gets flash message
func GetFlash(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return "", err
	}
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}

func encode(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func decode(src string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
