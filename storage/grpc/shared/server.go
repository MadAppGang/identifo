package shared

import (
	"encoding/json"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl model.UserStorage
	proto.UnimplementedUserStorageServer
}

func (m *GRPCServer) UserByPhone(ctx context.Context, in *proto.UserByPhoneRequest) (*proto.User, error) {
	user, err := m.Impl.UserByPhone(in.Phone)
	if err == model.ErrUserNotFound {
		return toProto(user), status.Error(codes.NotFound, err.Error())
	}
	return toProto(user), err
}

func (m *GRPCServer) AddUserWithPassword(ctx context.Context, in *proto.AddUserWithPasswordRequest) (*proto.User, error) {
	user, err := m.Impl.AddUserWithPassword(toModel(in.User), in.Password, in.Role, in.IsAnonymous)
	return toProto(user), err
}

func (m *GRPCServer) UserByID(ctx context.Context, in *proto.UserByIDRequest) (*proto.User, error) {
	user, err := m.Impl.UserByID(in.Id)
	if err == model.ErrUserNotFound {
		return toProto(user), status.Error(codes.NotFound, err.Error())
	}
	return toProto(user), err
}

func (m *GRPCServer) UserByEmail(ctx context.Context, in *proto.UserByEmailRequest) (*proto.User, error) {
	user, err := m.Impl.UserByEmail(in.Email)
	if err == model.ErrUserNotFound {
		return toProto(user), status.Error(codes.NotFound, err.Error())
	}
	return toProto(user), err
}

func (m *GRPCServer) UserByUsername(ctx context.Context, in *proto.UserByUsernameRequest) (*proto.User, error) {
	user, err := m.Impl.UserByUsername(in.Username)
	if err == model.ErrUserNotFound {
		return toProto(user), status.Error(codes.NotFound, err.Error())
	}
	return toProto(user), err
}

func (m *GRPCServer) UserByFederatedID(ctx context.Context, in *proto.UserByFederatedIDRequest) (*proto.User, error) {
	user, err := m.Impl.UserByFederatedID(in.Provider, in.Id)
	if err == model.ErrUserNotFound {
		return toProto(user), status.Error(codes.NotFound, err.Error())
	}
	return toProto(user), err
}

func (m *GRPCServer) AddUserWithFederatedID(ctx context.Context, in *proto.AddUserWithFederatedIDRequest) (*proto.User, error) {
	user, err := m.Impl.AddUserWithFederatedID(toModel(in.User), in.Provider, in.Id, in.Role)
	return toProto(user), err
}

func (m *GRPCServer) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.User, error) {
	user, err := m.Impl.UpdateUser(in.Id, toModel(in.User))
	return toProto(user), err
}

func (m *GRPCServer) ResetPassword(ctx context.Context, in *proto.ResetPasswordRequest) (*proto.Empty, error) {
	err := m.Impl.ResetPassword(in.Id, in.Password)
	return &proto.Empty{}, err
}

func (m *GRPCServer) CheckPassword(ctx context.Context, in *proto.CheckPasswordRequest) (*proto.Empty, error) {
	err := m.Impl.CheckPassword(in.Id, in.Password)
	return &proto.Empty{}, err
}

func (m *GRPCServer) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.Empty, error) {
	err := m.Impl.DeleteUser(in.Id)
	return &proto.Empty{}, err
}

func (m *GRPCServer) FetchUsers(ctx context.Context, in *proto.FetchUsersRequest) (*proto.FetchUsersResponse, error) {
	users, total, err := m.Impl.FetchUsers(in.Search, int(in.Skip), int(in.Limit))
	if err != nil {
		return &proto.FetchUsersResponse{}, err
	}

	protoUsers := []*proto.User{}
	for _, user := range users {
		protoUsers = append(protoUsers, toProto(user))
	}

	return &proto.FetchUsersResponse{
		Users:  protoUsers,
		Length: int32(total),
	}, nil
}

func (m *GRPCServer) UpdateLoginMetadata(ctx context.Context, in *proto.UpdateLoginMetadataRequest) (*proto.Empty, error) {
	payload := map[string]any{}

	_ = json.Unmarshal(in.PayloadJson, &payload)

	m.Impl.UpdateLoginMetadata(
		in.Operation,
		in.App,
		in.Id,
		in.Scopes,
		payload,
	)
	return &proto.Empty{}, nil
}

func (m *GRPCServer) AttachDeviceToken(ctx context.Context, in *proto.AttachDeviceTokenRequest) (*proto.Empty, error) {
	err := m.Impl.AttachDeviceToken(in.Id, in.Token)
	return &proto.Empty{}, err
}

func (m *GRPCServer) DetachDeviceToken(ctx context.Context, in *proto.DetachDeviceTokenRequest) (*proto.Empty, error) {
	err := m.Impl.DetachDeviceToken(in.Token)
	return &proto.Empty{}, err
}

func (m *GRPCServer) AllDeviceTokens(ctx context.Context, in *proto.AllDeviceTokensRequest) (*proto.AllDeviceTokensResponse, error) {
	return &proto.AllDeviceTokensResponse{}, nil
}

func (m *GRPCServer) ImportJSON(ctx context.Context, in *proto.ImportJSONRequest) (*proto.Empty, error) {
	err := m.Impl.ImportJSON(in.Data, in.ClearOldData)
	return &proto.Empty{}, err
}

func (m *GRPCServer) Close(ctx context.Context, in *proto.CloseRequest) (*proto.Empty, error) {
	m.Impl.Close()
	return &proto.Empty{}, nil
}
