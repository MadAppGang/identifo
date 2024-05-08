package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	pp "github.com/madappgang/identifo/v2/impersonation/grpc/impersonation_provider"
	grpcShared "github.com/madappgang/identifo/v2/impersonation/grpc/shared"
	"github.com/madappgang/identifo/v2/model"
	"google.golang.org/grpc"
)

const PluginName = "impersonation-provider"

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginName: &ImpersonationProviderPlugin{},
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type ImpersonationProviderPlugin struct {
	plugin.Plugin
	Impl model.ImpersonationProvider
}

func (p *ImpersonationProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pp.RegisterImpersonationProviderServer(s, &grpcShared.GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ImpersonationProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcShared.GRPCClient{Client: pp.NewImpersonationProviderClient(c)}, nil
}
