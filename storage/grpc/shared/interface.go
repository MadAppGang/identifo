package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/grpc/proto"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"user-storage": &UserStoragePlugin{},
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type UserStoragePlugin struct {
	plugin.Plugin
	Impl model.UserStorage
}

func (p *UserStoragePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterUserStorageServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *UserStoragePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewUserStorageClient(c)}, nil
}
