package model

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

// SessionService manages sessions.
type SessionService interface {
	NewSession() (Session, error)
	SessionDurationSeconds() int
	ProlongSession(sessionID string) error
}

// SessionManager is a default session service.
type SessionManager struct {
	sessionDuration SessionDuration
	sessionStorage  SessionStorage
}

// NewSessionManager creates new session manager and returns it.
func NewSessionManager(sessionDuration SessionDuration, sessionStorage SessionStorage) SessionService {
	return &SessionManager{
		sessionDuration: sessionDuration,
		sessionStorage:  sessionStorage,
	}
}

// NewSession creates new session and returns it.
func (sm *SessionManager) NewSession() (Session, error) {
	session := Session{ExpirationTime: time.Now().Add(sm.sessionDuration.Duration).Unix()}

	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return session, err
	}
	session.ID = base64.URLEncoding.EncodeToString(id)

	return session, nil
}

// SessionDurationSeconds returns session duration in seconds.
func (sm *SessionManager) SessionDurationSeconds() int {
	return int(sm.sessionDuration.Duration / time.Second)
}

// ProlongSession prolongs session duration.
func (sm *SessionManager) ProlongSession(sessionID string) error {
	err := sm.sessionStorage.ProlongSession(sessionID, sm.sessionDuration)
	return err
}
