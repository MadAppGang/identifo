package api

import (
	"errors"

	"github.com/madappgang/identifo/facebook"
)

// ErrFacebookEmptyUserID is when Facebook user ID is empty.
var ErrFacebookEmptyUserID = errors.New("Facebook user id is not accessible. ")

// FacebookUserID returns Facebook user ID.
func (ar *Router) FacebookUserID(accessToken string) (string, error) {
	fb := facebook.NewClient(accessToken)
	fbProfile, err := fb.MyProfile()
	if err != nil {
		return "", err
	}

	if len(fbProfile.ID) == 0 {
		return "", ErrFacebookEmptyUserID
	}
	return fbProfile.ID, nil
}
