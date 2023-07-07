package storage

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/madappgang/identifo/v2/l"
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
	zr := model.UserAuthChallenge{}
	app, err := c.as.AppByID(challenge.AppID)
	if err != nil {
		return zr, err
	}

	appAuthStrategies := app.AuthStrategies
	compatibleStrategies := model.FilterCompatible(challenge.Strategy, appAuthStrategies)
	// the app does not supports that type of challenge
	if len(compatibleStrategies) == 0 {
		return zr, l.LocalizedError{ErrID: l.ErrorRequestChallengeUnsupportedByAPP}
	}

	// selecting the first strategy from the list.
	// if there are more than one strategy we need to choose better one.
	auth := compatibleStrategies[0]

	// using the challenge he requested
	// if no user found, just silently return with no error for security reason
	u, err := c.UserByAuthStrategy(ctx, auth)
	if err != nil {
		return zr, nil
	}

	cha := challenge
	cha.UserID = u.ID
	cha.Strategy = auth
	cha.Solved = false
	cha.CreatedAt = time.Now()
	cha.ExpiresAt = cha.CreatedAt.Add(time.Minute * model.ExpireChallengeDuration(auth)) // challenge should know how long ti expire

	// this challenge type is about random OTP code, so generate it
	if auth.Type() == model.AuthStrategyFirstFactorInternal {
		f, ok := auth.(model.FirstFactorInternalStrategy)
		if ok {
			if f.Challenge == model.AuthChallengeTypeOTP ||
				f.Challenge == model.AuthChallengeTypeMagicLink {
				cha.OTP = randomOTP(6)
			}
		}
	}

	ch, err := c.uas.AddChallenge(ctx, cha)
	if err != nil {
		return zr, err
	}

	_ = c.sendChallengeToUser(ctx, cha, u)
	return ch, nil
}

func (c *UserStorageController) sendChallengeToUser(ctx context.Context, challenge model.UserAuthChallenge, u model.User) error {
	// no we need send the challenge to the user:
	// - sms
	// - email
	// - push: not supported yet
	// - socket: not supported yet
	// sending the challenge itself with the rendered text, each transport has it's own template
	// - magic link - build a link like reset password, but with OTP code to validate it
	// - OTP -- generate random number
	// using user
	return nil
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func randomOTP(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
