package controller

import "github.com/madappgang/identifo/v2/model"

func ephemeralUserForStrategy(strategy model.AuthStrategy, userIDValue string) model.User {
	result := model.User{
		ID:       model.NewUserID.String(),
		Username: model.NewUserUsername,
	}
	s, ok := strategy.(model.FirstFactorInternalStrategy)
	if ok {
		if s.Identity == model.AuthIdentityTypeEmail {
			result.Email = userIDValue
		} else if s.Identity == model.AuthIdentityTypePhone {
			result.PhoneNumber = userIDValue
		} else if s.Identity == model.AuthIdentityTypeUsername {
			result.Username = userIDValue
		} else if s.Identity == model.AuthIdentityTypeID {
			result.ID = userIDValue
		}
	}
	return result
}
