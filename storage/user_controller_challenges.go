package storage

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

// we need to get app, and check if it has strategies for this type of challenge
// if not - return it
// if yes - find the better (now just first suitable) strategy for it
// create full challenge
// send it to the user
// save it in db
// return it back
func (c *UserStorageController) RequestChallenge(ctx context.Context, challenge model.UserAuthChallenge) (model.UserAuthChallenge, error) {
	result := model.UserAuthChallenge{}
	app, err := c.as.AppByID(challenge.AppID)
	if err != nil {
		return result, err
	}

	appAuthStrategies := app.AuthStrategies
}
