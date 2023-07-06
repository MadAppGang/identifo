package shared

// import (
// 	"github.com/madappgang/identifo/v2/l"
// 	"github.com/madappgang/identifo/v2/model"
// 	"github.com/madappgang/identifo/v2/storage/grpc/proto"
// 	"golang.org/x/net/context"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// // GRPCClient is an implementation of KV that talks over RPC.
// type GRPCClient struct {
// 	Client         proto.UserStorageClient
// 	ClosableClient ClosableClient
// }

// func (m GRPCClient) UserByPhone(phone string) (model.User, error) {
// 	u, err := m.Client.UserByPhone(context.Background(), &proto.UserByPhoneRequest{
// 		Phone: phone,
// 	})
// 	if err != nil {
// 		if status.Convert(err).Code() == codes.NotFound {
// 			return model.User{}, l.ErrorUserNotFound
// 		}
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
// 	u, err := m.Client.AddUserWithPassword(context.Background(), &proto.AddUserWithPasswordRequest{
// 		User:        toProto(user),
// 		Password:    password,
// 		Role:        role,
// 		IsAnonymous: isAnonymous,
// 	})
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) UserByID(id string) (model.User, error) {
// 	u, err := m.Client.UserByID(context.Background(), &proto.UserByIDRequest{
// 		Id: id,
// 	})
// 	if err != nil {
// 		if status.Convert(err).Code() == codes.NotFound {
// 			return model.User{}, l.ErrorUserNotFound
// 		}
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) UserByEmail(email string) (model.User, error) {
// 	u, err := m.Client.UserByEmail(context.Background(), &proto.UserByEmailRequest{
// 		Email: email,
// 	})
// 	if err != nil {
// 		if status.Convert(err).Code() == codes.NotFound {
// 			return model.User{}, l.ErrorUserNotFound
// 		}
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) UserByUsername(username string) (model.User, error) {
// 	u, err := m.Client.UserByUsername(context.Background(), &proto.UserByUsernameRequest{
// 		Username: username,
// 	})
// 	if err != nil {
// 		if status.Convert(err).Code() == codes.NotFound {
// 			return model.User{}, l.ErrorUserNotFound
// 		}
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) UserByFederatedID(provider string, id string) (model.User, error) {
// 	u, err := m.Client.UserByFederatedID(context.Background(), &proto.UserByFederatedIDRequest{
// 		Id:       id,
// 		Provider: provider,
// 	})
// 	if err != nil {
// 		if status.Convert(err).Code() == codes.NotFound {
// 			return model.User{}, l.ErrorUserNotFound
// 		}
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) AddUserWithFederatedID(user model.User, provider string, id, role string) (model.User, error) {
// 	u, err := m.Client.AddUserWithFederatedID(context.Background(), &proto.AddUserWithFederatedIDRequest{
// 		User:     toProto(user),
// 		Provider: provider,
// 		Id:       id,
// 		Role:     role,
// 	})
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) UpdateUser(userID string, newUser model.User) (model.User, error) {
// 	u, err := m.Client.UpdateUser(context.Background(), &proto.UpdateUserRequest{
// 		User: toProto(newUser),
// 		Id:   userID,
// 	})
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return toModel(u), nil
// }

// func (m GRPCClient) ResetPassword(id, password string) error {
// 	_, err := m.Client.ResetPassword(context.Background(), &proto.ResetPasswordRequest{
// 		Id:       id,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m GRPCClient) CheckPassword(id, password string) error {
// 	_, err := m.Client.CheckPassword(context.Background(), &proto.CheckPasswordRequest{
// 		Id:       id,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m GRPCClient) DeleteUser(id string) error {
// 	_, err := m.Client.DeleteUser(context.Background(), &proto.DeleteUserRequest{
// 		Id: id,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m GRPCClient) FetchUsers(search string, skip, limit int) ([]model.User, int, error) {
// 	r, err := m.Client.FetchUsers(context.Background(), &proto.FetchUsersRequest{
// 		Search: search,
// 		Skip:   int32(skip),
// 		Limit:  int32(limit),
// 	})
// 	if err != nil {
// 		return []model.User{}, 0, err
// 	}

// 	users := []model.User{}

// 	for _, user := range r.Users {
// 		users = append(users, toModel(user))
// 	}

// 	return users, len(users), nil
// }

// func (m GRPCClient) UpdateLoginMetadata(userID string) {
// 	m.Client.UpdateLoginMetadata(context.Background(), &proto.UpdateLoginMetadataRequest{
// 		Id: userID,
// 	})
// }

// // push device tokens
// func (m GRPCClient) AttachDeviceToken(userID, token string) error {
// 	_, err := m.Client.AttachDeviceToken(context.Background(), &proto.AttachDeviceTokenRequest{
// 		Id:    userID,
// 		Token: token,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m GRPCClient) DetachDeviceToken(token string) error {
// 	_, err := m.Client.DetachDeviceToken(context.Background(), &proto.DetachDeviceTokenRequest{
// 		Token: token,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m GRPCClient) AllDeviceTokens(userID string) ([]string, error) {
// 	r, err := m.Client.AllDeviceTokens(context.Background(), &proto.AllDeviceTokensRequest{
// 		Id: userID,
// 	})
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	return r.Tokens, nil
// }

// // import data
// func (m GRPCClient) ImportJSON(data []byte, clearOldData bool) error {
// 	return nil
// }

// func (m GRPCClient) Close() {
// 	m.Client.Close(context.Background(), &proto.CloseRequest{})
// 	if m.ClosableClient != nil {
// 		m.ClosableClient.Close()
// 	}
// }

// func toModel(u *proto.User) model.User {
// 	return model.User{
// 		ID:       u.Id,
// 		Username: u.Username,
// 		Email:    u.Email,
// 		FullName: u.FullName,
// 		Phone:    u.Phone,
// 		Pswd:     u.Pswd,
// 		Active:   u.Active,
// 		TFAInfo: model.TFAInfo{
// 			IsEnabled:     u.TfaInfo.IsEnabled,
// 			HOTPCounter:   int(u.TfaInfo.HotpCounter),
// 			HOTPExpiredAt: u.TfaInfo.HotpExpiredAt.AsTime(),
// 			Secret:        u.TfaInfo.Secret,
// 		},
// 		NumOfLogins:     int(u.NumOfLogins),
// 		LatestLoginTime: u.LatestLoginTime,
// 		AccessRole:      u.AccessRole,
// 		Anonymous:       u.Anonymous,
// 		FederatedIDs:    u.FederatedIds,
// 		Scopes:          u.Scopes,
// 	}
// }

// func toProto(u model.User) *proto.User {
// 	return &proto.User{
// 		Id:       u.ID,
// 		Username: u.Username,
// 		Email:    u.Email,
// 		FullName: u.FullName,
// 		Phone:    u.Phone,
// 		Pswd:     u.Pswd,
// 		Active:   u.Active,
// 		TfaInfo: &proto.User_TFAInfo{
// 			IsEnabled:     u.TFAInfo.IsEnabled,
// 			HotpCounter:   int32(u.TFAInfo.HOTPCounter),
// 			HotpExpiredAt: timestamppb.New(u.TFAInfo.HOTPExpiredAt),
// 			Secret:        u.TFAInfo.Secret,
// 		},
// 		NumOfLogins:     int32(u.NumOfLogins),
// 		LatestLoginTime: u.LatestLoginTime,
// 		AccessRole:      u.AccessRole,
// 		Anonymous:       u.Anonymous,
// 		FederatedIds:    u.FederatedIDs,
// 		Scopes:          u.Scopes,
// 	}
// }

// // Here is the gRPC server that GRPCClient talks to.
// type GRPCServer struct {
// 	// This is the real implementation
// 	Impl model.UserStorage
// 	proto.UnimplementedUserStorageServer
// }

// func (m *GRPCServer) UserByPhone(ctx context.Context, in *proto.UserByPhoneRequest) (*proto.User, error) {
// 	user, err := m.Impl.UserByPhone(in.Phone)
// 	if err == l.ErrorUserNotFound {
// 		return toProto(user), status.Errorf(codes.NotFound, err.Error())
// 	}
// 	return toProto(user), err
// }

// func (m *GRPCServer) AddUserWithPassword(ctx context.Context, in *proto.AddUserWithPasswordRequest) (*proto.User, error) {
// 	user, err := m.Impl.AddUserWithPassword(toModel(in.User), in.Password, in.Role, in.IsAnonymous)
// 	return toProto(user), err
// }

// func (m *GRPCServer) UserByID(ctx context.Context, in *proto.UserByIDRequest) (*proto.User, error) {
// 	user, err := m.Impl.UserByID(in.Id)
// 	if err == l.ErrorUserNotFound {
// 		return toProto(user), status.Errorf(codes.NotFound, err.Error())
// 	}
// 	return toProto(user), err
// }

// func (m *GRPCServer) UserByEmail(ctx context.Context, in *proto.UserByEmailRequest) (*proto.User, error) {
// 	user, err := m.Impl.UserByEmail(in.Email)
// 	if err == l.ErrorUserNotFound {
// 		return toProto(user), status.Errorf(codes.NotFound, err.Error())
// 	}
// 	return toProto(user), err
// }

// func (m *GRPCServer) UserByUsername(ctx context.Context, in *proto.UserByUsernameRequest) (*proto.User, error) {
// 	user, err := m.Impl.UserByUsername(in.Username)
// 	if err == l.ErrorUserNotFound {
// 		return toProto(user), status.Errorf(codes.NotFound, err.Error())
// 	}
// 	return toProto(user), err
// }

// func (m *GRPCServer) UserByFederatedID(ctx context.Context, in *proto.UserByFederatedIDRequest) (*proto.User, error) {
// 	user, err := m.Impl.UserByFederatedID(in.Provider, in.Id)
// 	if err == l.ErrorUserNotFound {
// 		return toProto(user), status.Errorf(codes.NotFound, err.Error())
// 	}
// 	return toProto(user), err
// }

// func (m *GRPCServer) AddUserWithFederatedID(ctx context.Context, in *proto.AddUserWithFederatedIDRequest) (*proto.User, error) {
// 	user, err := m.Impl.AddUserWithFederatedID(toModel(in.User), in.Provider, in.Id, in.Role)
// 	return toProto(user), err
// }

// func (m *GRPCServer) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.User, error) {
// 	user, err := m.Impl.UpdateUser(in.Id, toModel(in.User))
// 	return toProto(user), err
// }

// func (m *GRPCServer) ResetPassword(ctx context.Context, in *proto.ResetPasswordRequest) (*proto.Empty, error) {
// 	err := m.Impl.ResetPassword(in.Id, in.Password)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) CheckPassword(ctx context.Context, in *proto.CheckPasswordRequest) (*proto.Empty, error) {
// 	err := m.Impl.CheckPassword(in.Id, in.Password)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.Empty, error) {
// 	err := m.Impl.DeleteUser(in.Id)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) FetchUsers(ctx context.Context, in *proto.FetchUsersRequest) (*proto.FetchUsersResponse, error) {
// 	users, total, err := m.Impl.FetchUsers(in.Search, int(in.Skip), int(in.Limit))
// 	if err != nil {
// 		return &proto.FetchUsersResponse{}, err
// 	}

// 	protoUsers := []*proto.User{}
// 	for _, user := range users {
// 		protoUsers = append(protoUsers, toProto(user))
// 	}

// 	return &proto.FetchUsersResponse{
// 		Users:  protoUsers,
// 		Length: int32(total),
// 	}, nil
// }

// func (m *GRPCServer) UpdateLoginMetadata(ctx context.Context, in *proto.UpdateLoginMetadataRequest) (*proto.Empty, error) {
// 	m.Impl.UpdateLoginMetadata(in.Id)
// 	return &proto.Empty{}, nil
// }

// func (m *GRPCServer) AttachDeviceToken(ctx context.Context, in *proto.AttachDeviceTokenRequest) (*proto.Empty, error) {
// 	err := m.Impl.AttachDeviceToken(in.Id, in.Token)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) DetachDeviceToken(ctx context.Context, in *proto.DetachDeviceTokenRequest) (*proto.Empty, error) {
// 	err := m.Impl.DetachDeviceToken(in.Token)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) AllDeviceTokens(ctx context.Context, in *proto.AllDeviceTokensRequest) (*proto.AllDeviceTokensResponse, error) {
// 	return &proto.AllDeviceTokensResponse{}, nil
// }

// func (m *GRPCServer) ImportJSON(ctx context.Context, in *proto.ImportJSONRequest) (*proto.Empty, error) {
// 	err := m.Impl.ImportJSON(in.Data, in.ClearOldData)
// 	return &proto.Empty{}, err
// }

// func (m *GRPCServer) Close(ctx context.Context, in *proto.CloseRequest) (*proto.Empty, error) {
// 	m.Impl.Close()
// 	return &proto.Empty{}, nil
// }
