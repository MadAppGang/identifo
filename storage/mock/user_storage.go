package mock

import (
	"context"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

type UserStorage struct {
	Storage
	Users []model.User
	UData map[string]model.UserData
}

func (us *UserStorage) UserByID(ctx context.Context, id string) (model.User, error) {
	for _, u := range us.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, l.ErrorUserNotFound
}

func (us *UserStorage) UserBySecondaryID(ctx context.Context, idt model.AuthIdentityType, id string) (model.User, error) {
	for _, u := range us.Users {
		if idt == model.AuthIdentityTypeEmail && u.Email == id {
			return u, nil
		} else if idt == model.AuthIdentityTypePhone && u.PhoneNumber == id {
			return u, nil
		}
	}
	return model.User{}, l.ErrorUserNotFound
}

func (us *UserStorage) UserByFederatedID(ctx context.Context, idType model.UserFederatedType, userIdentityTypeOther, externalID string) (model.User, error) {
	return model.User{}, l.ErrorLoginTypeNotSupported
}

func (us *UserStorage) UserData(ctx context.Context, userID string, fields ...model.UserDataField) (model.UserData, error) {
	d, ok := us.UData[userID]
	if !ok {
		return model.UserData{}, l.ErrorUserNotFound
	}

	return d, nil
}

func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	return nil
}
