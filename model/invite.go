package model

import (
	"errors"
	"time"
)

// Invite is a representation of the invite model.
// Token field is required for proper working.
type Invite struct {
	ID        string    `json:"id" bson:"_id"`
	AppID     string    `json:"appId" bson:"appId"`
	Token     string    `json:"token" bson:"token"`
	Valid     bool      `json:"valid" bson:"valid"`
	Email     string    `json:"email" bson:"email"`
	Role      string    `json:"role" bson:"role"`
	CreatedBy string    `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}

func (i Invite) Validate() error {
	if i.Email == "" {
		return errors.New("email cannot be empty")
	}
	if i.ExpiresAt.IsZero() {
		return errors.New("expiresAt cannot represents the zero time instant")
	}
	return nil
}
