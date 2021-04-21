package dynamodb

import (
	"time"

	"github.com/madappgang/identifo/model"
)

const invitesTableName = "Invites"

// InviteStorage is a DynamoDB invite storage.
type InviteStorage struct {
	db *DB
}

// NewInviteStorage creates new DynamoDB invite storage.
func NewInviteStorage(db *DB) (model.InviteStorage, error) {
	is := &InviteStorage{db: db}
	err := is.ensureTable()
	return is, err
}

// ensureTable ensures that token storage exists in the database.
func (is *InviteStorage) ensureTable() error {
	panic("implement me")
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	panic("implement me")
}

// GetByEmail returns valid and not expired invite by email.
func (is *InviteStorage) GetByEmail(email string) (model.Invite, error) {
	panic("implement me")
}

// GetByID returns invite by its ID.
func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	panic("implement me")
}

// GetAll returns all active invites by default.
// To get an invalid invites need to set withInvalid argument to true.
func (is *InviteStorage) GetAll(withInvalid bool, skip, limit int) ([]model.Invite, int, error) {
	panic("implement me")
}

// InvalidateAllByEmail invalidates all invites by email.
func (is *InviteStorage) InvalidateAllByEmail(email string) error {
	panic("implement me")
}

// InvalidateByID invalidates specific invite by its ID.
func (is *InviteStorage) InvalidateByID(id string) error {
	panic("implement me")
}
