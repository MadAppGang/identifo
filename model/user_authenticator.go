package model

import "time"

// UserAuthEnrolment user enrollment to one of auth strategies.
type UserAuthEnrolment struct {
	ID          string        `json:"id,omitempty"`
	UserID      string        `json:"user_id,omitempty"`
	Name        string        `json:"name,omitempty"`
	StrategyID  string        `json:"strategy_id,omitempty"`
	Strategy    *AuthStrategy `json:"strategy,omitempty"`
	Confirmed   bool          `json:"confirmed,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	EnrolledAt  time.Time     `json:"enrolled_at,omitempty"`
	ConfirmedAt time.Time     `json:"confirmed_at,omitempty"`
	LastAuthAt  time.Time     `json:"last_auth_at,omitempty"`
}
