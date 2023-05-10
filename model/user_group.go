package model

import "time"

// Group is a group related to tenant
type Group struct {
	ID        string    `json:"id,omitempty"`
	Default   bool      `json:"default,omitempty"`
	Name      string    `json:"name,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// some default groups
const (
	GroupDefault   = "default"
	GroupAdmin     = "admin"
	GroupModerator = "moderator"
)
