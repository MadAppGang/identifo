package mem

import (
	"time"

	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

// InviteStorage is an in-memory invite storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type InviteStorage struct {
	storage map[string]model.Invite
}

// NewInviteStorage creates an in-memory invite storage.
func NewInviteStorage() (model.InviteStorage, error) {
	return &InviteStorage{storage: make(map[string]model.Invite)}, nil
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	is.storage[inviteToken] = model.Invite{
		ID:        xid.New().String(),
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
	for _, invite := range is.storage {
		if invite.Email == email && invite.Valid == true && invite.ExpiresAt.After(time.Now()) {
			return invite, nil
		}
	}
	return model.Invite{}, model.ErrorNotFound
}

// GetByID returns invite by its ID.
func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	invite, ok := is.storage[id]
	if !ok {
		return model.Invite{}, model.ErrorNotFound
	}
	return invite, nil
}

// GetAll returns all active invites by default.
// To get an invalid invites need to set withInvalid argument to true.
func (is *InviteStorage) GetAll(withInvalid bool, skip, limit int) ([]model.Invite, int, error) {
	var (
		invites []model.Invite
		total   int
	)

	for _, invite := range is.storage {
		if !withInvalid && invite.Valid {
			total++
			skip--
			if skip > -1 || (limit != 0 && len(invites) == limit) {
				break
			}
			invites = append(invites, invite)
		}
	}

	return invites, total, nil
}

// InvalidateAllByEmail invalidates all invites by email.
func (is *InviteStorage) InvalidateAllByEmail(email string) error {
	for _, invite := range is.storage {
		if invite.Email == email {
			invite.Valid = false
			is.storage[invite.ID] = invite
		}
	}
	return nil
}

// InvalidateByID invalidates specific invite by its ID.
func (is *InviteStorage) InvalidateByID(id string) error {
	invite, ok := is.storage[id]
	if !ok {
		return model.ErrorNotFound
	}

	invite.Valid = false
	is.storage[invite.ID] = invite
	return nil
}
