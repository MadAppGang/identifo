package mem

import (
	"time"

	"github.com/madappgang/identifo/model"
)

// TokenStorage is an in-memory token storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type InviteStorage struct {
	storage map[string]model.Invite
}

// NewTokenStorage creates an in-memory token storage.
func NewInviteStorage() (model.InviteStorage, error) {
	return &InviteStorage{storage: make(map[string]model.Invite)}, nil
}

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
