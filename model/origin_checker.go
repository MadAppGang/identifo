package model

import "net/http"

type OriginChecker interface {
	With(f func(r *http.Request, origin string) bool) OriginChecker
	CheckOrigin(r *http.Request, origin string) bool
	AddRawURLs(urls []string)
	DeleteAll()
}
