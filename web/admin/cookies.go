package admin

import (
	"encoding/base64"
)

const (
	cookieName = "SessionID"
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
