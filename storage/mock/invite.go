package mock

import (
	"context"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/exp/maps"
)

var inviteNotFoundError = l.NewError(l.ErrorNotFound, "invite")

// InviteStorage is an in-memory invite storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type InviteStorage struct {
	Storage
	Invites map[string]model.Invite
}

func NewInviteStorage() *InviteStorage {
	return &InviteStorage{Invites: make(map[string]model.Invite)}
}

func (s *InviteStorage) Save(ctx context.Context, invite model.Invite) error {
	s.Invites[invite.ID] = invite
	return nil
}

func (s *InviteStorage) GetByID(ctx context.Context, id string) (model.Invite, error) {
	i, ok := s.Invites[id]
	if !ok {
		return model.Invite{}, inviteNotFoundError
	}
	return i, nil
}

func (s *InviteStorage) GetAll(ctx context.Context, withArchived bool, skip, limit int) ([]model.Invite, int, error) {
	return maps.Values(s.Invites), len(s.Invites), nil
}

func (s *InviteStorage) Update(ctx context.Context, invite model.Invite) error {
	s.Invites[invite.ID] = invite
	return nil
}
