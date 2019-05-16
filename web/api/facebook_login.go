package api

import (
	"errors"

	"github.com/madappgang/identifo/facebook"
)

var ErrFacebookEmptyUserID = errors.New("Facebook user id is not accessible. ")

// FacebookLogin implements federated facebook login
func (ar *Router) FacebookUserID(accessToken string) (string, error) {
	fb := facebook.NewClient(accessToken)
	fbProfile, err := fb.MyProfile()
	if err != nil {
		return "", err
	}

	// check we had `id` permissions for the access_token
	if len(fbProfile.ID) == 0 {
		return "", ErrFacebookEmptyUserID
	}
	return fbProfile.ID, nil
}
