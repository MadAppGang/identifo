package model

import (
	"time"
)

// SessionStorage is an interface for session storage.
type SessionStorage interface {
	GetSession(id string) (Session, error)
	InsertSession(session Session) error
	DeleteSession(id string) error
	ProlongSession(id string, newDuration time.Duration) error
}
