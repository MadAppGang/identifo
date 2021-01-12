// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/proto"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "IDENTIFO_USER_STORAGE_PLUGIN",
	MagicCookieValue: "ON",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"user_storage": &UserStorageGRPCPlugin{},
}

// UserStorage is the interface that we're exposing as a plugin.
type UserStorage interface {
	UserByPhone(phone string) (*proto.User, error)
	AddUserByPhone(phone, role string) (*proto.User, error)
	UserByID(id string) (*proto.User, error)
	UserByEmail(email string) (*proto.User, error)
	IDByName(name string) (string, error)
	AttachDeviceToken(id, token string) error
	DetachDeviceToken(token string) error
	UserByNamePassword(name, password string) (*proto.User, error)
	AddUserByNameAndPassword(username, password, role string, isAnonymous bool) (*proto.User, error)
	UserExists(name string) bool
	UserByFederatedID(provider proto.FederatedIdentityProvider, id string) (*proto.User, error)
	AddUserWithFederatedID(provider proto.FederatedIdentityProvider, id, role string) (*proto.User, error)
	UpdateUser(userID string, newUser *proto.User) (*proto.User, error)
	ResetPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]*proto.User, int, error)

	RequestScopes(userID string, scopes []string) ([]string, error)
	Scopes() []string
	ImportJSON(data []byte) error
	UpdateLoginMetadata(userID string)
	Close()
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type UserStorageGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl UserStorage
}

func (p *UserStorageGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterUserStorageServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *UserStorageGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewUserStorageClient(c)}, nil
}
