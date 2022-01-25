package model

import (
	"time"
)

// InviteStorage is a storage for invites.
type InviteStorage interface {
	Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error
	GetByEmail(email string) (Invite, error)
	GetByID(id string) (Invite, error)
	GetAll(withArchived bool, skip, limit int) ([]Invite, int, error)
	ArchiveAllByEmail(email string) error
	ArchiveByID(id string) error
	Close()
}
