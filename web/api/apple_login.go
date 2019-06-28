package api

import (
	"errors"

	"github.com/madappgang/identifo/identity_providers/apple"
	"github.com/madappgang/identifo/model"
)

// ErrAppleEmptyUserID is when Apple user ID is empty.
var ErrAppleEmptyUserID = errors.New("Apple user id is not accessible. ")

// AppleUserID returns Apple user ID.
func (ar *Router) AppleUserID(authorizationCode string, appleInfo *model.AppleInfo) (string, error) {
	ac := apple.NewClient(authorizationCode, appleInfo)
	appleProfile, err := ac.MyProfile()
	if err != nil {
		return "", err
	}

	if len(appleProfile.ID) == 0 {
		return "", ErrAppleEmptyUserID
	}
	return appleProfile.ID, nil
}
