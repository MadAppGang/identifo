package mem

import (
	"time"

	"github.com/madappgang/identifo/model"
)

// TokenStorage is an in-memory invite storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type InviteStorage struct {
	storage map[string]model.Invite
}

// NewInviteStorage creates an in-memory invite storage.
func NewInviteStorage() (model.InviteStorage, error) {
	return &InviteStorage{storage: make(map[string]model.Invite)}, nil
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(inviteToken, email, role, appID, createdBy string, expiresAt time.Time) error {
	is.storage[inviteToken] = model.Invite{
		AppID:     appID,
		Token:     inviteToken,
		Valid:     true,
		Email:     email,
		Role:      role,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	return nil
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
