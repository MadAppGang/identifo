package mem

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

func (ms *ManagementKeysStorage) CreateKey(ctx context.Context, name string, scopes []string) (model.ManagementKey, error) {
}

func (ms *ManagementKeysStorage) DisableKey(ctx context.Context, id string) (model.ManagementKey, error) {
}

func (ms *ManagementKeysStorage) RenameKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
}

func (ms *ManagementKeysStorage) ChangeScopesForKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
}

func (ms *ManagementKeysStorage) UseKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
}

func (ms *ManagementKeysStorage) GeyAllKeys(ctx context.Context) ([]model.ManagementKey, error) {
}
