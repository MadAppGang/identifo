package shared

import (
	"io"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	Client   proto.UserStorageClient
	Closable io.Closer
}

func (m GRPCClient) UserByPhone(phone string) (model.User, error) {
	u, err := m.Client.UserByPhone(context.Background(), &proto.UserByPhoneRequest{
		Phone: phone,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	u, err := m.Client.AddUserWithPassword(context.Background(), &proto.AddUserWithPasswordRequest{
		User:        toProto(user),
		Password:    password,
		Role:        role,
		IsAnonymous: isAnonymous,
	})
	if err != nil {
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) UserByID(id string) (model.User, error) {
	u, err := m.Client.UserByID(context.Background(), &proto.UserByIDRequest{
		Id: id,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) UserByEmail(email string) (model.User, error) {
	u, err := m.Client.UserByEmail(context.Background(), &proto.UserByEmailRequest{
		Email: email,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) UserByUsername(username string) (model.User, error) {
	u, err := m.Client.UserByUsername(context.Background(), &proto.UserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) UserByFederatedID(provider string, id string) (model.User, error) {
	u, err := m.Client.UserByFederatedID(context.Background(), &proto.UserByFederatedIDRequest{
		Id:       id,
		Provider: provider,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) AddUserWithFederatedID(user model.User, provider string, id, role string) (model.User, error) {
	u, err := m.Client.AddUserWithFederatedID(context.Background(), &proto.AddUserWithFederatedIDRequest{
		User:     toProto(user),
		Provider: provider,
		Id:       id,
		Role:     role,
	})
	if err != nil {
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) UpdateUser(userID string, newUser model.User) (model.User, error) {
	u, err := m.Client.UpdateUser(context.Background(), &proto.UpdateUserRequest{
		User: toProto(newUser),
		Id:   userID,
	})
	if err != nil {
		return model.User{}, err
	}

	return toModel(u), nil
}

func (m GRPCClient) ResetPassword(id, password string) error {
	_, err := m.Client.ResetPassword(context.Background(), &proto.ResetPasswordRequest{
		Id:       id,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m GRPCClient) CheckPassword(id, password string) error {
	_, err := m.Client.CheckPassword(context.Background(), &proto.CheckPasswordRequest{
		Id:       id,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m GRPCClient) DeleteUser(id string) error {
	_, err := m.Client.DeleteUser(context.Background(), &proto.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m GRPCClient) FetchUsers(search string, skip, limit int) ([]model.User, int, error) {
	r, err := m.Client.FetchUsers(context.Background(), &proto.FetchUsersRequest{
		Search: search,
		Skip:   int32(skip),
		Limit:  int32(limit),
	})
	if err != nil {
		return []model.User{}, 0, err
	}

	users := []model.User{}

	for _, user := range r.Users {
		users = append(users, toModel(user))
	}

	return users, len(users), nil
}

func (m GRPCClient) UpdateLoginMetadata(userID string) {
	m.Client.UpdateLoginMetadata(context.Background(), &proto.UpdateLoginMetadataRequest{
		Id: userID,
	})
}

// push device tokens
func (m GRPCClient) AttachDeviceToken(userID, token string) error {
	_, err := m.Client.AttachDeviceToken(context.Background(), &proto.AttachDeviceTokenRequest{
		Id:    userID,
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m GRPCClient) DetachDeviceToken(token string) error {
	_, err := m.Client.DetachDeviceToken(context.Background(), &proto.DetachDeviceTokenRequest{
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m GRPCClient) AllDeviceTokens(userID string) ([]string, error) {
	r, err := m.Client.AllDeviceTokens(context.Background(), &proto.AllDeviceTokensRequest{
		Id: userID,
	})
	if err != nil {
		return []string{}, err
	}

	return r.Tokens, nil
}

// import data
func (m GRPCClient) ImportJSON(data []byte, clearOldData bool) error {
	return nil
}

func (m GRPCClient) Close() {
	m.Client.Close(context.Background(), &proto.CloseRequest{})
	if m.Closable != nil {
		m.Closable.Close()
	}
}

func toModel(u *proto.User) model.User {
	return model.User{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		FullName: u.FullName,
		Phone:    u.Phone,
		Pswd:     u.Pswd,
		Active:   u.Active,
		TFAInfo: model.TFAInfo{
			IsEnabled:     u.TfaInfo.IsEnabled,
			HOTPCounter:   int(u.TfaInfo.HotpCounter),
			HOTPExpiredAt: u.TfaInfo.HotpExpiredAt.AsTime(),
			Secret:        u.TfaInfo.Secret,
		},
		NumOfLogins:     int(u.NumOfLogins),
		LatestLoginTime: u.LatestLoginTime,
		AccessRole:      u.AccessRole,
		Anonymous:       u.Anonymous,
		FederatedIDs:    u.FederatedIds,
		Scopes:          u.Scopes,
	}
}

func toProto(u model.User) *proto.User {
	return &proto.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		FullName: u.FullName,
		Phone:    u.Phone,
		Pswd:     u.Pswd,
		Active:   u.Active,
		TfaInfo: &proto.User_TFAInfo{
			IsEnabled:     u.TFAInfo.IsEnabled,
			HotpCounter:   int32(u.TFAInfo.HOTPCounter),
			HotpExpiredAt: timestamppb.New(u.TFAInfo.HOTPExpiredAt),
			Secret:        u.TFAInfo.Secret,
		},
		NumOfLogins:     int32(u.NumOfLogins),
		LatestLoginTime: u.LatestLoginTime,
		AccessRole:      u.AccessRole,
		Anonymous:       u.Anonymous,
		FederatedIds:    u.FederatedIDs,
		Scopes:          u.Scopes,
	}
}
