package mock

import (
	"context"
	"errors"
	"time"

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
	if len(challenge.ID) == 0 {
		challenge.ID = model.NewID().String()
	}
	u.Challenges = append(u.Challenges, challenge)
	return challenge, nil
}

func (u *UserAuthStorage) GetLatestChallenge(ctx context.Context, strategy model.AuthStrategy, userID string) (model.UserAuthChallenge, error) {
	// just return the last one
	return u.Challenges[len(u.Challenges)-1], nil
}

func (u *UserAuthStorage) MarkChallengeAsSolved(ctx context.Context, challenge model.UserAuthChallenge) error {
	for i, ch := range u.Challenges {
		if ch.ID == challenge.ID {
			challenge.Solved = true
			challenge.SolvedAt = time.Now()
			u.Challenges[i] = challenge
			return nil
		}
	}
	return errors.New("not found")
}
