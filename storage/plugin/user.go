package plugin

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	grpcShared "github.com/madappgang/identifo/v2/storage/grpc/shared"
	"github.com/madappgang/identifo/v2/storage/plugin/shared"
)

// NewUserStorage creates and inits plugin user storage.
func NewUserStorage(settings model.PluginSettings) (model.UserStorage, error) {
	var err error
	params := []string{}
	for k, v := range settings.Params {
		params = append(params, "-"+k)
		params = append(params, v)
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command(settings.Cmd, params...),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginName)
	if err != nil {
		return nil, err
	}

	user := raw.(*grpcShared.GRPCClient)

	user.Closable = pluginClosableClient{client: client}

	return user, err
}

type pluginClosableClient struct {
	client *plugin.Client
}

func (g pluginClosableClient) Close() error {
	g.client.Kill()
	return nil
}
