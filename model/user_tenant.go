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

// some default roles
const (
	RoleOwner = "owner"
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

// TenantMembership is representation for user tenant membership
// tenant have a list of groups and
// and user can have multiply roles in group
type TenantMembership struct {
	TenantID   string              `json:"tenant_id,omitempty"`
	TenantName string              `json:"tenant_name,omitempty"`
	Groups     map[string][]string `json:"groups,omitempty"`
}
