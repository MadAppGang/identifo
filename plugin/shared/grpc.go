package shared

import (
	"context"
	"errors"

	"github.com/madappgang/identifo/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.UserStorageClient }

func (m *GRPCClient) UserByPhone(phone string) (*proto.User, error) {
	user, err := m.client.UserByPhone(context.Background(), &proto.UserByPhoneRequest{
		Phone: phone,
	})
	return user, err
}

func (m *GRPCClient) UserByEmail(email string) (*proto.User, error) {
	user, err := m.client.UserByEmail(context.Background(), &proto.UserByEmailRequest{
		Email: email,
	})
	return user, err
}

func (m *GRPCClient) UserByFederatedID(provider proto.FederatedIdentityProvider, fid string) (*proto.User, error) {
	user, err := m.client.UserByFederatedID(context.Background(), &proto.UserByFederatedIDRequest{
		Id:       fid,
		Provider: provider,
	})
	return user, err
}

func (m *GRPCClient) UserByID(id string) (*proto.User, error) {
	user, err := m.client.UserByID(context.Background(), &proto.UserByIDRequest{
		Id: id,
	})
	return user, err
}

func (m *GRPCClient) UserByNamePassword(name, password string) (*proto.User, error) {
	user, err := m.client.UserByNamePassword(context.Background(), &proto.UserByNamePasswordRequest{
		Name:     name,
		Password: password,
	})
	return user, err
}

func (m *GRPCClient) UserExists(name string) bool {
	res, err := m.client.UserExists(context.Background(), &proto.UserExistsRequest{
		Name: name,
	})
	if err != nil {
		return false
	}
	return res.DoesExist
}

func (m *GRPCClient) AddUserByNameAndPassword(name, password, role string, isAnonymous bool) (*proto.User, error) {
	user, err := m.client.AddUserByNameAndPassword(context.Background(), &proto.AddUserByNameAndPasswordRequest{
		Name:        name,
		Password:    password,
		Role:        role,
		IsAnonymous: isAnonymous,
	})
	return user, err
}

func (m *GRPCClient) AddUserByPhone(phone, role string) (*proto.User, error) {
	user, err := m.client.AddUserByPhone(context.Background(), &proto.AddUserByPhoneRequest{
		Phone: phone,
		Role:  role,
	})
	return user, err
}

func (m *GRPCClient) AddUserWithFederatedID(provider proto.FederatedIdentityProvider, id, role string) (*proto.User, error) {
	user, err := m.client.AddUserWithFederatedID(context.Background(), &proto.AddUserWithFederatedIDRequest{
		Provider: provider,
		Id:       id,
		Role:     role,
	})
	return user, err
}

func (m *GRPCClient) AttachDeviceToken(id, token string) error {
	_, err := m.client.AttachDeviceToken(context.Background(), &proto.AttachDeviceTokenRequest{
		Id:    id,
		Token: token,
	})
	return err
}

func (m *GRPCClient) DetachDeviceToken(token string) error {
	_, err := m.client.DetachDeviceToken(context.Background(), &proto.DetachDeviceTokenRequest{
		Token: token,
	})
	return err
}

func (m *GRPCClient) Close() {
	m.client.Close(context.Background(), new(proto.Empty))
}

func (m *GRPCClient) DeleteUser(id string) error {
	_, err := m.client.DeleteUser(context.Background(), &proto.DeleteUserRequest{
		Id: id,
	})
	return err
}

func (m *GRPCClient) FetchUsers(search string, skip, limit int) ([]*proto.User, int, error) {
	res, err := m.client.FetchUsers(context.Background(), &proto.FetchUsersRequest{
		Search: search,
		Skip:   uint32(skip),
		Limit:  uint32(limit),
	})
	if err != nil {
		return nil, 0, err
	}
	return res.Users, int(res.Total), err
}

func (m *GRPCClient) IDByName(name string) (string, error) {
	res, err := m.client.IDByName(context.Background(), &proto.IDByNameRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return res.Id, err
}

func (m *GRPCClient) ImportJSON(data []byte) error {
	_, err := m.client.ImportJSON(context.Background(), &proto.ImportJSONRequest{
		Data: data,
	})
	return err
}

func (m *GRPCClient) RequestScopes(userID string, scopes []string) ([]string, error) {
	res, err := m.client.RequestScopes(context.Background(), &proto.RequestScopesRequest{
		UserId: userID,
		Scopes: scopes,
	})
	if err != nil {
		return nil, err
	}
	return res.Scopes, err
}

func (m *GRPCClient) Scopes() []string {
	res, err := m.client.Scopes(context.Background(), &proto.Empty{})
	if err != nil {
		return nil
	}
	return res.Scopes
}

func (m *GRPCClient) ResetPassword(userID, password string) error {
	_, err := m.client.ResetPassword(context.Background(), &proto.ResetPasswordRequest{
		Id:       userID,
		Password: password,
	})
	return err
}

func (m *GRPCClient) UpdateLoginMetadata(userID string) {
	m.client.UpdateLoginMetadata(context.Background(), &proto.UpdateLoginMetadataRequest{
		UserId: userID,
	})
}

func (m *GRPCClient) UpdateUser(userID string, newUser *proto.User) (*proto.User, error) {
	user, err := m.client.UpdateUser(context.Background(), &proto.UpdateUserRequest{
		UserId:  userID,
		NewUser: newUser,
	})
	return user, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl UserStorage
}

func (m *GRPCServer) UserByID(ctx context.Context, req *proto.UserByIDRequest) (*proto.User, error) {
	return m.Impl.UserByID(req.Id)
}

func (m *GRPCServer) UserByPhone(ctx context.Context, req *proto.UserByPhoneRequest) (*proto.User, error) {
	return m.Impl.UserByPhone(req.Phone)
}

func (m *GRPCServer) UserByEmail(ctx context.Context, req *proto.UserByEmailRequest) (*proto.User, error) {
	return m.Impl.UserByEmail(req.Email)
}

func (m *GRPCServer) UserExists(ctx context.Context, req *proto.UserExistsRequest) (*proto.UserExistsResponse, error) {
	return &proto.UserExistsResponse{
		DoesExist: m.Impl.UserExists(req.Name),
	}, nil
}

func (m *GRPCServer) UserByNamePassword(ctx context.Context, req *proto.UserByNamePasswordRequest) (*proto.User, error) {
	return m.Impl.UserByNamePassword(req.Name, req.Password)
}

func (m *GRPCServer) UserByFederatedID(ctx context.Context, req *proto.UserByFederatedIDRequest) (*proto.User, error) {
	return m.Impl.UserByFederatedID(req.Provider, req.Id)
}

func (m *GRPCServer) AddUserByNameAndPassword(ctx context.Context, req *proto.AddUserByNameAndPasswordRequest) (*proto.User, error) {
	return m.Impl.AddUserByNameAndPassword(req.Name, req.Password, req.Role, req.IsAnonymous)
}

func (m *GRPCServer) AddUserByPhone(ctx context.Context, req *proto.AddUserByPhoneRequest) (*proto.User, error) {
	return m.Impl.AddUserByPhone(req.Phone, req.Role)
}

func (m *GRPCServer) AddUserWithFederatedID(ctx context.Context, req *proto.AddUserWithFederatedIDRequest) (*proto.User, error) {
	return m.Impl.AddUserWithFederatedID(req.Provider, req.Id, req.Role)
}

func (m *GRPCServer) AttachDeviceToken(ctx context.Context, req *proto.AttachDeviceTokenRequest) (*proto.Empty, error) {
	return new(proto.Empty), m.Impl.AttachDeviceToken(req.Id, req.Token)
}

func (m *GRPCServer) DetachDeviceToken(ctx context.Context, req *proto.DetachDeviceTokenRequest) (*proto.Empty, error) {
	return new(proto.Empty), m.Impl.DetachDeviceToken(req.Token)
}

func (m *GRPCServer) Close(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	m.Impl.Close()
	return new(proto.Empty), nil
}

func (m *GRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.Empty, error) {
	return new(proto.Empty), m.Impl.DeleteUser(req.Id)
}

func (m *GRPCServer) FetchUsers(ctx context.Context, req *proto.FetchUsersRequest) (*proto.FetchUsersResponse, error) {
	users, total, err := m.Impl.FetchUsers(req.Search, int(req.Skip), int(req.Limit))
	return &proto.FetchUsersResponse{
		Users: users,
		Total: uint32(total),
	}, err
}

func (m *GRPCServer) IDByName(ctx context.Context, req *proto.IDByNameRequest) (*proto.IDByNameResponse, error) {
	id, err := m.Impl.IDByName(req.Name)
	return &proto.IDByNameResponse{
		Id: id,
	}, err
}

func (m *GRPCServer) ImportJSON(ctx context.Context, req *proto.ImportJSONRequest) (*proto.Empty, error) {
	return new(proto.Empty), m.Impl.ImportJSON(req.Data)
}

func (m *GRPCServer) RequestScopes(ctx context.Context, req *proto.RequestScopesRequest) (*proto.ScopesResponse, error) {
	scopes, err := m.Impl.RequestScopes(req.UserId, req.Scopes)
	return &proto.ScopesResponse{
		Scopes: scopes,
	}, err
}

func (m *GRPCServer) Scopes(ctx context.Context, req *proto.Empty) (*proto.ScopesResponse, error) {
	return &proto.ScopesResponse{
		Scopes: m.Impl.Scopes(),
	}, nil
}

func (m *GRPCServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.Empty, error) {
	err := m.Impl.ResetPassword(req.Id, req.Password)
	return &proto.Empty{}, err
}

func (m *GRPCServer) UpdateLoginMetadata(ctx context.Context, req *proto.UpdateLoginMetadataRequest) (*proto.Empty, error) {
	m.Impl.UpdateLoginMetadata(req.UserId)
	return &proto.Empty{}, nil
}

func (m *GRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
	if req.NewUser == nil {
		return nil, errors.New("User is nil")
	}
	user, err := m.Impl.UpdateUser(req.UserId, req.NewUser)
	return user, err
}
