package model

import "time"

// Group is a group related to tenant
type Group struct {
	ID        string
	Default   bool
	Name      string
	Tags      []string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// some default groups
const (
	GroupDefault   = "default"
	GroupAdmin     = "admin"
	GroupModerator = "moderator"
)
