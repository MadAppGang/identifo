package mem

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/exp/maps"
)

type ManagementKeysStorage struct {
	storage map[string]model.ManagementKey
}

// NewKeysManagementStorage creates an in-memory management keys storage.
func NewManagementKeysStorage() (model.ManagementKeysStorage, error) {
	return &ManagementKeysStorage{storage: make(map[string]model.ManagementKey)}, nil
}

func (ms *ManagementKeysStorage) GetKey(ctx context.Context, id string) (model.ManagementKey, error) {
	key, ok := ms.storage[id]
	if !ok {
		return key, errors.New("not found")
	}
	return key, nil
}

func (ms *ManagementKeysStorage) CreateKey(ctx context.Context, name string, scopes []string) (model.ManagementKey, error) {
	key := model.ManagementKey{
		Name:      name,
		Scopes:    scopes,
		ID:        uuid.New().String(),
		Active:    true,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	ms.storage[key.ID] = key
	return key, nil
}

func (ms *ManagementKeysStorage) DisableKey(ctx context.Context, id string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Active = false
	ms.storage[key.ID] = key
	return key, nil
}

func (ms *ManagementKeysStorage) RenameKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Name = name
	ms.storage[key.ID] = key
	return key, nil
}

func (ms *ManagementKeysStorage) ChangeScopesForKey(ctx context.Context, id string, scopes []string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Scopes = scopes
	ms.storage[key.ID] = key
	return key, nil
}

func (ms *ManagementKeysStorage) UseKey(ctx context.Context, id string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.LastUsed = time.Now()
	ms.storage[key.ID] = key
	return key, nil
}

func (ms *ManagementKeysStorage) GeyAllKeys(ctx context.Context) ([]model.ManagementKey, error) {
	keys := maps.Values(ms.storage)
	return keys, nil
}

func (ms *ManagementKeysStorage) ImportJSON(data []byte, cleanOldData bool) error {
	return errors.New("not implemented")
}
