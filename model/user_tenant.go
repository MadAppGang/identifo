package model

import "time"

const TenantDefault = "default"

// Tenant is an abstract way to user exists in some isolated company area.
type Tenant struct {
	ID        string
	Name      string
	Default   bool
	Tags      []string
	UpdatedAt time.Time
	CreatedAt time.Time
}
