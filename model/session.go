package model

import (
	"time"
)

// Session is a session.
type Session struct {
	ID             string
	ExpirationDate time.Time
}
