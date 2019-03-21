package model

import (
	"time"
)

// Session is a session.
type Session struct {
	ID             string
	ExpirationDate time.Time
}

// SessionDuration wraps time.Duration to implement custom yaml and json encoding and decoding.
type SessionDuration struct {
	time.Duration
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (sd *SessionDuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if sd == nil {
		return nil
	}

	var aux int
	if err := unmarshal(&aux); err != nil {
		return err
	}

	*sd = SessionDuration{Duration: time.Second * time.Duration(aux)}
	return nil
}
