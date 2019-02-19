package html

import (
	"encoding/base64"
	"net/http"
	"time"
)

const (
	// CookieKeyUserID cookie key to keep authenticated user's ID.
	CookieKeyUserID = "identifo-user"
)

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

func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	c := &http.Cookie{Name: name, Value: encode(value), MaxAge: maxAge, HttpOnly: true}
	http.SetCookie(w, c)
}

func getCookie(r *http.Request, name string) (string, error) {
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

	return value, nil
}

func deleteCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{Name: name, Value: "", Expires: time.Unix(0, 0), MaxAge: -1}
	http.SetCookie(w, c)
}
