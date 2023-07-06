package storage

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

// TODO: implement challenge creation
func (c *UserStorageController) RequestChallenge(ctx context.Context, challenge model.UserAuthChallenge) (model.UserAuthChallenge, error) {
}
