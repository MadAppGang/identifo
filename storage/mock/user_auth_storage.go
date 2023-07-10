package mock

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

type UserAuthStorage struct {
	Storage
	Challenges []model.UserAuthChallenge
}

func (u *UserAuthStorage) ImportJSON(data []byte, clearOldData bool) error {
	return nil
}

func (u *UserAuthStorage) AddChallenge(ctx context.Context, challenge model.UserAuthChallenge) (model.UserAuthChallenge, error) {
	u.Challenges = append(u.Challenges, challenge)
	return challenge, nil
}

func (u *UserAuthStorage) GetLatestChallenge(ctx context.Context, strategy model.AuthStrategy, userID string) (model.UserAuthChallenge, error) {
	// just return the last one
	return u.Challenges[len(u.Challenges)-1], nil
}
