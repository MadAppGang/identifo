package model

import "time"

const TenantDefault = "default"

// Tenant is an abstract way to user exists in some isolated company area.
type Tenant struct {
	ID        string
	Name      string
	Default   bool
	Tags      []string
	Groups    []Group
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Group is a group related to tenant
type Group struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Default   bool      `json:"default,omitempty"`
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
