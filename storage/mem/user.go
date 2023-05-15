package mem

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/model"
)

// NewUserStorage creates and inits in-memory user storage.
// Use it only for test purposes and in CI, all data is wiped on exit.
func NewUserStorage() (*UserStorage, error) {
	return &UserStorage{
		users:       []model.User{},
		userData:    make(map[string]model.UserData),
		userDevices: make(map[string]string),
	}, nil
}

// UserStorage is an in-memory user storage .
type UserStorage struct {
	users       []model.User
	userData    map[string]model.UserData
	userDevices map[string]string
}

// ================================================================
// UserStorage implementations
// ================================================================

// UserByID returns user from memory storage.
func (us *UserStorage) UserByID(ctx context.Context, ID string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.ID, ID) {
			return u, nil
		}
	}
	return model.User{}, model.ErrUserNotFound
}

// UserByPhone returns user from memory storage by phone.
func (us *UserStorage) UserByPhone(ctx context.Context, Phone string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.PhoneNumber, Phone) {
			return u, nil
		}
	}
	return model.User{}, model.ErrUserNotFound
}

// UserByPhone returns user from memory storage by email.
func (us *UserStorage) UserByEmail(ctx context.Context, Email string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.Email, Email) {
			return u, nil
		}
	}
	return model.User{}, model.ErrUserNotFound
}

// UserByPhone returns user from memory storage by identity.
func (us *UserStorage) UserByIdentity(ctx context.Context, idType model.UserIdentityType, userIdentityTypeOther, externalID string) (model.User, error) {
	for _, u := range us.userData {
		for _, i := range u.Identities {
			if i.Type == idType && i.ExternalID == externalID && i.TypeOther == userIdentityTypeOther {
				return us.UserByID(ctx, u.UserID)
			}
		}
	}
	return model.User{}, model.ErrUserNotFound
}

// UserData returns user data for user for specific fields.
func (us *UserStorage) UserData(ctx context.Context, userID string, fields ...model.UserDataField) (model.UserData, error) {
	data, ok := us.userData[userID]
	if !ok {
		return model.UserData{}, model.ErrUserNotFound
	}

	result := model.FilterUserDataFields(data, fields...)
	return result, nil
}

func (us *UserStorage) AddUser(ctx context.Context, user model.User) (model.User, error) {
	if len(user.ID) == 0 {
		user.ID = fmt.Sprintf("id:%d", time.Now().UnixNano())
	}
	us.users = append(us.users, user)
	return user, nil
}

func (us *UserStorage) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	for i, u := range us.users {
		if strings.EqualFold(u.ID, user.ID) {
			us.users[i] = user
			return user, nil
		}
	}
	return user, model.ErrUserNotFound
}

func (us *UserStorage) UpdateUserData(ctx context.Context, userID string, data model.UserData, fields ...model.UserDataField) (model.UserData, error) {
	d, ok := us.userData[userID]
	if !ok {
		return data, model.ErrUserNotFound
	}

	for _, f := range fields {
		switch f {
		case model.UserDataFieldTenantMembership:
			d.TenantMembership = data.TenantMembership
		case model.UserDataFieldAuthEnrollments:
			d.AuthEnrollments = data.AuthEnrollments
		case model.UserDataFieldIdentities:
			d.Identities = data.Identities
		case model.UserDataFieldMFAEnrollments:
			d.MFAEnrollments = data.MFAEnrollments
		case model.UserDataFieldActiveDevices:
			d.ActiveDevices = data.ActiveDevices
		case model.UserDataFieldAppsData:
			d.AppsData = data.AppsData
		case model.UserDataFieldData:
			d.Data = data.Data
		case model.UserDataFieldAll:
			d = data
		default:
		}
	}

	us.userData[userID] = d
	return d, nil
}

// ================================================================
// Storage implementations
// ================================================================
func (us *UserStorage) Ready(ctx context.Context) error {
	return nil
}

func (us *UserStorage) Connect(ctx context.Context) error {
	return nil
}

func (us *UserStorage) Close(ctx context.Context) error {
	return nil
}

// ================================================================
// ImportableStorage implementations
// ================================================================
func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.users = []model.User{}
		us.userData = make(map[string]model.UserData)
		us.userDevices = make(map[string]string)
	}

	ud := model.UserImportData{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	us.users = ud.Users
	for _, u := range ud.Data {
		us.userData[u.UserID] = u
	}
	return nil
}
