package controller

import (
	"context"
	"errors"
	"time"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// RefreshJWTToken issues new access and, if requested, refresh token for provided refresh token.
// After new tokens are issued, the old refresh token and access token gets invalidated (added to blocklist).
// We validate refresh token
// if its valid - issue new tokens.
// ! Be careful, old access token still could be accepted by some systems, if it is not yet expired and those systems are not checking blocklist (usually the should not in distributed systems).
func (c *UserStorageController) RefreshJWTToken(ctx context.Context, refresh_token *model.JWToken, access string, app model.AppData, scopes []string) (model.AuthResponse, error) {
	sub, err := refresh_token.Claims.GetSubject()
	if err != nil || len(sub) == 0 {
		return model.AuthResponse{}, l.ErrorValidationTokenInvalidSubject
	}

	u, err := c.u.UserByID(ctx, sub)
	if err != nil {
		return model.AuthResponse{}, err
	}

	// let's check if the refresh token is blocked
	_, err = c.toks.TokenByID(ctx, refresh_token.ID())
	if err != nil && !errors.Is(err, l.ErrorNotFound) {
		return model.AuthResponse{}, nil
	} else if err == nil {
		return model.AuthResponse{}, l.ErrorTokenBlocked
	}

	// let's parse access token:
	at, err := jwt.ParseTokenString(access)
	if err == nil {
		tse := model.TokenStorageEntity{
			ID:        at.ID(),
			RawToken:  access,
			TokenType: model.TokenTypeAccess,
			AddedAt:   time.Now(),
			AddedBy:   model.TokenStorageAddedByUser,
			Comments:  "Refresh token API call",
		}
		c.toks.SaveToken(ctx, tse)
	}
	tse := model.TokenStorageEntity{
		ID:        refresh_token.ID(),
		RawToken:  refresh_token.Raw,
		TokenType: model.TokenTypeRefresh,
		AddedAt:   time.Now(),
		AddedBy:   model.TokenStorageAddedByUser,
		Comments:  "Refresh token API call",
	}

	err = c.toks.SaveToken(ctx, tse)
	if err != nil {
		return model.AuthResponse{}, err
	}

	response, err := c.GetJWTTokens(ctx, app, u, scopes)
	if err != nil {
		return model.AuthResponse{}, err
	}
	return response, nil
}