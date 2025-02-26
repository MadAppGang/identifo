// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.2
// source: storage/grpc/proto/user.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	UserStorage_UserByPhone_FullMethodName            = "/proto.UserStorage/UserByPhone"
	UserStorage_AddUserWithPassword_FullMethodName    = "/proto.UserStorage/AddUserWithPassword"
	UserStorage_UserByID_FullMethodName               = "/proto.UserStorage/UserByID"
	UserStorage_UserByEmail_FullMethodName            = "/proto.UserStorage/UserByEmail"
	UserStorage_UserByUsername_FullMethodName         = "/proto.UserStorage/UserByUsername"
	UserStorage_UserByFederatedID_FullMethodName      = "/proto.UserStorage/UserByFederatedID"
	UserStorage_AddUserWithFederatedID_FullMethodName = "/proto.UserStorage/AddUserWithFederatedID"
	UserStorage_UpdateUser_FullMethodName             = "/proto.UserStorage/UpdateUser"
	UserStorage_ResetPassword_FullMethodName          = "/proto.UserStorage/ResetPassword"
	UserStorage_CheckPassword_FullMethodName          = "/proto.UserStorage/CheckPassword"
	UserStorage_DeleteUser_FullMethodName             = "/proto.UserStorage/DeleteUser"
	UserStorage_FetchUsers_FullMethodName             = "/proto.UserStorage/FetchUsers"
	UserStorage_UpdateLoginMetadata_FullMethodName    = "/proto.UserStorage/UpdateLoginMetadata"
	UserStorage_AttachDeviceToken_FullMethodName      = "/proto.UserStorage/AttachDeviceToken"
	UserStorage_DetachDeviceToken_FullMethodName      = "/proto.UserStorage/DetachDeviceToken"
	UserStorage_AllDeviceTokens_FullMethodName        = "/proto.UserStorage/AllDeviceTokens"
	UserStorage_ImportJSON_FullMethodName             = "/proto.UserStorage/ImportJSON"
	UserStorage_Close_FullMethodName                  = "/proto.UserStorage/Close"
)

// UserStorageClient is the client API for UserStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserStorageClient interface {
	// UserByPhone(phone string) (model.User, error) {
	UserByPhone(ctx context.Context, in *UserByPhoneRequest, opts ...grpc.CallOption) (*User, error)
	// AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	AddUserWithPassword(ctx context.Context, in *AddUserWithPasswordRequest, opts ...grpc.CallOption) (*User, error)
	// UserByID(id string) (model.User, error) {
	UserByID(ctx context.Context, in *UserByIDRequest, opts ...grpc.CallOption) (*User, error)
	// UserByEmail(email string) (model.User, error) {
	UserByEmail(ctx context.Context, in *UserByEmailRequest, opts ...grpc.CallOption) (*User, error)
	// UserByUsername(username string) (model.User, error) {
	UserByUsername(ctx context.Context, in *UserByUsernameRequest, opts ...grpc.CallOption) (*User, error)
	// UserByFederatedID(provider string, id string) (model.User, error) {
	UserByFederatedID(ctx context.Context, in *UserByFederatedIDRequest, opts ...grpc.CallOption) (*User, error)
	// AddUserWithFederatedID(user model.User, provider string, id, role string) (model.User, error) {
	AddUserWithFederatedID(ctx context.Context, in *AddUserWithFederatedIDRequest, opts ...grpc.CallOption) (*User, error)
	// UpdateUser(userID string, newUser model.User) (model.User, error) {
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error)
	// ResetPassword(id, password string) error {
	ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*Empty, error)
	// CheckPassword(id, password string) error {
	CheckPassword(ctx context.Context, in *CheckPasswordRequest, opts ...grpc.CallOption) (*Empty, error)
	// DeleteUser(id string) error {
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*Empty, error)
	// FetchUsers(search string, skip, limit int) ([]model.User, int, error) {
	FetchUsers(ctx context.Context, in *FetchUsersRequest, opts ...grpc.CallOption) (*FetchUsersResponse, error)
	// UpdateLoginMetadata(userID string) {
	UpdateLoginMetadata(ctx context.Context, in *UpdateLoginMetadataRequest, opts ...grpc.CallOption) (*Empty, error)
	// push device tokens
	// AttachDeviceToken(userID, token string) error {
	AttachDeviceToken(ctx context.Context, in *AttachDeviceTokenRequest, opts ...grpc.CallOption) (*Empty, error)
	// DetachDeviceToken(token string) error {
	DetachDeviceToken(ctx context.Context, in *DetachDeviceTokenRequest, opts ...grpc.CallOption) (*Empty, error)
	// AllDeviceTokens(userID string) ([]string, error) {
	AllDeviceTokens(ctx context.Context, in *AllDeviceTokensRequest, opts ...grpc.CallOption) (*AllDeviceTokensResponse, error)
	// import data
	// ImportJSON(data []byte) error {
	ImportJSON(ctx context.Context, in *ImportJSONRequest, opts ...grpc.CallOption) (*Empty, error)
	// Close() {
	Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*Empty, error)
}

type userStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewUserStorageClient(cc grpc.ClientConnInterface) UserStorageClient {
	return &userStorageClient{cc}
}

func (c *userStorageClient) UserByPhone(ctx context.Context, in *UserByPhoneRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UserByPhone_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) AddUserWithPassword(ctx context.Context, in *AddUserWithPasswordRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_AddUserWithPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UserByID(ctx context.Context, in *UserByIDRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UserByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UserByEmail(ctx context.Context, in *UserByEmailRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UserByEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UserByUsername(ctx context.Context, in *UserByUsernameRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UserByUsername_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UserByFederatedID(ctx context.Context, in *UserByFederatedIDRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UserByFederatedID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) AddUserWithFederatedID(ctx context.Context, in *AddUserWithFederatedIDRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_AddUserWithFederatedID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserStorage_UpdateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_ResetPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) CheckPassword(ctx context.Context, in *CheckPasswordRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_CheckPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_DeleteUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) FetchUsers(ctx context.Context, in *FetchUsersRequest, opts ...grpc.CallOption) (*FetchUsersResponse, error) {
	out := new(FetchUsersResponse)
	err := c.cc.Invoke(ctx, UserStorage_FetchUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) UpdateLoginMetadata(ctx context.Context, in *UpdateLoginMetadataRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_UpdateLoginMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) AttachDeviceToken(ctx context.Context, in *AttachDeviceTokenRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_AttachDeviceToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) DetachDeviceToken(ctx context.Context, in *DetachDeviceTokenRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_DetachDeviceToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) AllDeviceTokens(ctx context.Context, in *AllDeviceTokensRequest, opts ...grpc.CallOption) (*AllDeviceTokensResponse, error) {
	out := new(AllDeviceTokensResponse)
	err := c.cc.Invoke(ctx, UserStorage_AllDeviceTokens_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) ImportJSON(ctx context.Context, in *ImportJSONRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_ImportJSON_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userStorageClient) Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, UserStorage_Close_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserStorageServer is the server API for UserStorage service.
// All implementations must embed UnimplementedUserStorageServer
// for forward compatibility
type UserStorageServer interface {
	// UserByPhone(phone string) (model.User, error) {
	UserByPhone(context.Context, *UserByPhoneRequest) (*User, error)
	// AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	AddUserWithPassword(context.Context, *AddUserWithPasswordRequest) (*User, error)
	// UserByID(id string) (model.User, error) {
	UserByID(context.Context, *UserByIDRequest) (*User, error)
	// UserByEmail(email string) (model.User, error) {
	UserByEmail(context.Context, *UserByEmailRequest) (*User, error)
	// UserByUsername(username string) (model.User, error) {
	UserByUsername(context.Context, *UserByUsernameRequest) (*User, error)
	// UserByFederatedID(provider string, id string) (model.User, error) {
	UserByFederatedID(context.Context, *UserByFederatedIDRequest) (*User, error)
	// AddUserWithFederatedID(user model.User, provider string, id, role string) (model.User, error) {
	AddUserWithFederatedID(context.Context, *AddUserWithFederatedIDRequest) (*User, error)
	// UpdateUser(userID string, newUser model.User) (model.User, error) {
	UpdateUser(context.Context, *UpdateUserRequest) (*User, error)
	// ResetPassword(id, password string) error {
	ResetPassword(context.Context, *ResetPasswordRequest) (*Empty, error)
	// CheckPassword(id, password string) error {
	CheckPassword(context.Context, *CheckPasswordRequest) (*Empty, error)
	// DeleteUser(id string) error {
	DeleteUser(context.Context, *DeleteUserRequest) (*Empty, error)
	// FetchUsers(search string, skip, limit int) ([]model.User, int, error) {
	FetchUsers(context.Context, *FetchUsersRequest) (*FetchUsersResponse, error)
	// UpdateLoginMetadata(userID string) {
	UpdateLoginMetadata(context.Context, *UpdateLoginMetadataRequest) (*Empty, error)
	// push device tokens
	// AttachDeviceToken(userID, token string) error {
	AttachDeviceToken(context.Context, *AttachDeviceTokenRequest) (*Empty, error)
	// DetachDeviceToken(token string) error {
	DetachDeviceToken(context.Context, *DetachDeviceTokenRequest) (*Empty, error)
	// AllDeviceTokens(userID string) ([]string, error) {
	AllDeviceTokens(context.Context, *AllDeviceTokensRequest) (*AllDeviceTokensResponse, error)
	// import data
	// ImportJSON(data []byte) error {
	ImportJSON(context.Context, *ImportJSONRequest) (*Empty, error)
	// Close() {
	Close(context.Context, *CloseRequest) (*Empty, error)
	mustEmbedUnimplementedUserStorageServer()
}

// UnimplementedUserStorageServer must be embedded to have forward compatible implementations.
type UnimplementedUserStorageServer struct {
}

func (UnimplementedUserStorageServer) UserByPhone(context.Context, *UserByPhoneRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserByPhone not implemented")
}
func (UnimplementedUserStorageServer) AddUserWithPassword(context.Context, *AddUserWithPasswordRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserWithPassword not implemented")
}
func (UnimplementedUserStorageServer) UserByID(context.Context, *UserByIDRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserByID not implemented")
}
func (UnimplementedUserStorageServer) UserByEmail(context.Context, *UserByEmailRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserByEmail not implemented")
}
func (UnimplementedUserStorageServer) UserByUsername(context.Context, *UserByUsernameRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserByUsername not implemented")
}
func (UnimplementedUserStorageServer) UserByFederatedID(context.Context, *UserByFederatedIDRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserByFederatedID not implemented")
}
func (UnimplementedUserStorageServer) AddUserWithFederatedID(context.Context, *AddUserWithFederatedIDRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserWithFederatedID not implemented")
}
func (UnimplementedUserStorageServer) UpdateUser(context.Context, *UpdateUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserStorageServer) ResetPassword(context.Context, *ResetPasswordRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPassword not implemented")
}
func (UnimplementedUserStorageServer) CheckPassword(context.Context, *CheckPasswordRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPassword not implemented")
}
func (UnimplementedUserStorageServer) DeleteUser(context.Context, *DeleteUserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserStorageServer) FetchUsers(context.Context, *FetchUsersRequest) (*FetchUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchUsers not implemented")
}
func (UnimplementedUserStorageServer) UpdateLoginMetadata(context.Context, *UpdateLoginMetadataRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoginMetadata not implemented")
}
func (UnimplementedUserStorageServer) AttachDeviceToken(context.Context, *AttachDeviceTokenRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachDeviceToken not implemented")
}
func (UnimplementedUserStorageServer) DetachDeviceToken(context.Context, *DetachDeviceTokenRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachDeviceToken not implemented")
}
func (UnimplementedUserStorageServer) AllDeviceTokens(context.Context, *AllDeviceTokensRequest) (*AllDeviceTokensResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllDeviceTokens not implemented")
}
func (UnimplementedUserStorageServer) ImportJSON(context.Context, *ImportJSONRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportJSON not implemented")
}
func (UnimplementedUserStorageServer) Close(context.Context, *CloseRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedUserStorageServer) mustEmbedUnimplementedUserStorageServer() {}

// UnsafeUserStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserStorageServer will
// result in compilation errors.
type UnsafeUserStorageServer interface {
	mustEmbedUnimplementedUserStorageServer()
}

func RegisterUserStorageServer(s grpc.ServiceRegistrar, srv UserStorageServer) {
	s.RegisterService(&UserStorage_ServiceDesc, srv)
}

func _UserStorage_UserByPhone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByPhoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UserByPhone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UserByPhone_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UserByPhone(ctx, req.(*UserByPhoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_AddUserWithPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserWithPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).AddUserWithPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_AddUserWithPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).AddUserWithPassword(ctx, req.(*AddUserWithPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UserByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UserByID(ctx, req.(*UserByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UserByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UserByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UserByEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UserByEmail(ctx, req.(*UserByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UserByUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UserByUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UserByUsername_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UserByUsername(ctx, req.(*UserByUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UserByFederatedID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByFederatedIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UserByFederatedID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UserByFederatedID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UserByFederatedID(ctx, req.(*UserByFederatedIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_AddUserWithFederatedID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserWithFederatedIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).AddUserWithFederatedID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_AddUserWithFederatedID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).AddUserWithFederatedID(ctx, req.(*AddUserWithFederatedIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_ResetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).ResetPassword(ctx, req.(*ResetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_CheckPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).CheckPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_CheckPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).CheckPassword(ctx, req.(*CheckPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_FetchUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).FetchUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_FetchUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).FetchUsers(ctx, req.(*FetchUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_UpdateLoginMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLoginMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).UpdateLoginMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_UpdateLoginMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).UpdateLoginMetadata(ctx, req.(*UpdateLoginMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_AttachDeviceToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttachDeviceTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).AttachDeviceToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_AttachDeviceToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).AttachDeviceToken(ctx, req.(*AttachDeviceTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_DetachDeviceToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetachDeviceTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).DetachDeviceToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_DetachDeviceToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).DetachDeviceToken(ctx, req.(*DetachDeviceTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_AllDeviceTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllDeviceTokensRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).AllDeviceTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_AllDeviceTokens_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).AllDeviceTokens(ctx, req.(*AllDeviceTokensRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_ImportJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).ImportJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_ImportJSON_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).ImportJSON(ctx, req.(*ImportJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserStorage_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserStorageServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserStorage_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserStorageServer).Close(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserStorage_ServiceDesc is the grpc.ServiceDesc for UserStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserStorage",
	HandlerType: (*UserStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserByPhone",
			Handler:    _UserStorage_UserByPhone_Handler,
		},
		{
			MethodName: "AddUserWithPassword",
			Handler:    _UserStorage_AddUserWithPassword_Handler,
		},
		{
			MethodName: "UserByID",
			Handler:    _UserStorage_UserByID_Handler,
		},
		{
			MethodName: "UserByEmail",
			Handler:    _UserStorage_UserByEmail_Handler,
		},
		{
			MethodName: "UserByUsername",
			Handler:    _UserStorage_UserByUsername_Handler,
		},
		{
			MethodName: "UserByFederatedID",
			Handler:    _UserStorage_UserByFederatedID_Handler,
		},
		{
			MethodName: "AddUserWithFederatedID",
			Handler:    _UserStorage_AddUserWithFederatedID_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserStorage_UpdateUser_Handler,
		},
		{
			MethodName: "ResetPassword",
			Handler:    _UserStorage_ResetPassword_Handler,
		},
		{
			MethodName: "CheckPassword",
			Handler:    _UserStorage_CheckPassword_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserStorage_DeleteUser_Handler,
		},
		{
			MethodName: "FetchUsers",
			Handler:    _UserStorage_FetchUsers_Handler,
		},
		{
			MethodName: "UpdateLoginMetadata",
			Handler:    _UserStorage_UpdateLoginMetadata_Handler,
		},
		{
			MethodName: "AttachDeviceToken",
			Handler:    _UserStorage_AttachDeviceToken_Handler,
		},
		{
			MethodName: "DetachDeviceToken",
			Handler:    _UserStorage_DetachDeviceToken_Handler,
		},
		{
			MethodName: "AllDeviceTokens",
			Handler:    _UserStorage_AllDeviceTokens_Handler,
		},
		{
			MethodName: "ImportJSON",
			Handler:    _UserStorage_ImportJSON_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _UserStorage_Close_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage/grpc/proto/user.proto",
}
