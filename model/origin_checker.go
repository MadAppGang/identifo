package model

import "net/http"

type OriginCheckFunc func(r *http.Request, origin string) bool

type OriginChecker interface {
	AddCheck(f OriginCheckFunc)
	CheckOrigin(r *http.Request, origin string) bool
	Update() error
}
