package model

import (
	"context"
)

// InviteStorage is a storage for invites.
type InviteStorage interface {
	Storage
	ImportableStorage

	Save(ctx context.Context, invite Invite) error
	// GetByEmail(ctx context.Context, email, tenant, groups string) ([]Invite, error)
	// GetForTenant(ctx context.Context, tenant, group string) ([]Invite, error)
	// GetByID(ctx context.Context, id string) (Invite, error)
	GetAll(ctx context.Context, withArchived bool, skip, limit int) ([]Invite, int, error)

	Update(ctx context.Context, invite Invite) error
}
