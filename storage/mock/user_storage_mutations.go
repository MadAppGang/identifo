package mock

import (
	"context"
	"errors"

	"github.com/madappgang/identifo/v2/model"
)

// // User mutation
func (us *UserStorage) AddUser(ctx context.Context, user model.User) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (us *UserStorage) UpdateUser(ctx context.Context, user model.User, fields ...string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (us *UserStorage) UpdateUserData(ctx context.Context, userID string, data model.UserData, fields ...model.UserDataField) (model.UserData, error) {
	us.UData[userID] = data
	return data, nil
}

func (us *UserStorage) DeleteUser(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}
