package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	pp "github.com/madappgang/identifo/v2/user_payload_provider/grpc/payload_provider"
	grpcShared "github.com/madappgang/identifo/v2/user_payload_provider/grpc/shared"
	"google.golang.org/grpc"
)

const PluginName = "token-payload-provider"

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginName: &TokenPayloadProviderPlugin{},
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type TokenPayloadProviderPlugin struct {
	plugin.Plugin
	Impl model.TokenPayloadProvider
}

func (p *TokenPayloadProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pp.RegisterPayloadProviderServiceServer(s, &grpcShared.GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TokenPayloadProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcShared.GRPCClient{Client: pp.NewPayloadProviderServiceClient(c)}, nil
}
