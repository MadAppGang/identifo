package dynamodb

import (
	"time"

	"github.com/madappgang/identifo/model"
)

const invitesTableName = "Invites"

// TokenStorage is a DynamoDB token storage.
type InviteStorage struct {
	db *DB
}

// NewTokenStorage creates new DynamoDB token storage.
func NewInviteStorage(db *DB) (model.InviteStorage, error) {
	is := &InviteStorage{db: db}
	err := is.ensureTable()
	return is, err
}

// ensureTable ensures that token storage exists in the database.
func (is *InviteStorage) ensureTable() error {
	panic("implement me")
}

func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	panic("implement me")
}

func (is *InviteStorage) GetByEmail(email string) (model.Invite, error) {
	panic("implement me")
}

func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	panic("implement me")
}

func (is *InviteStorage) GetAll(withInvalid bool, skip, limit int) ([]model.Invite, int, error) {
	panic("implement me")
}

func (is *InviteStorage) InvalidateAllByEmail(email string) error {
	panic("implement me")
}

func (is *InviteStorage) InvalidateByID(id string) error {
	panic("implement me")
}
