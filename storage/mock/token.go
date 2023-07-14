package mock

import (
	"context"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

var tokenNotFoundError = l.NewError(l.ErrorNotFound, "token")

// NewTokenStorage creates an in-memory token storage.
func NewTokenStorage() (model.TokenStorage, error) {
	return &TokenStorage{storage: make(map[string]model.TokenStorageEntity)}, nil
}

// TokenStorage is an in-memory token storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type TokenStorage struct {
	storage map[string]model.TokenStorageEntity
}

func (s *TokenStorage) SaveToken(ctx context.Context, token model.TokenStorageEntity) error {
	s.storage[token.ID] = token
	return nil
}

func (s *TokenStorage) TokenByID(ctx context.Context, id string) (model.TokenStorageEntity, error) {
	t, ok := s.storage[id]
	if !ok {
		return model.TokenStorageEntity{}, tokenNotFoundError
	}
	return t, nil
}

func (s *TokenStorage) TokenByRaw(ctx context.Context, raw string) (model.TokenStorageEntity, error) {
	for _, t := range s.storage {
		if t.RawToken == raw {
			return t, nil
		}
	}
	return model.TokenStorageEntity{}, tokenNotFoundError
}

func (s *TokenStorage) DeleteToken(ctx context.Context, id string) error {
	delete(s.storage, id)
	return nil
}

func (s *TokenStorage) Close() {
	s.storage = nil
}
