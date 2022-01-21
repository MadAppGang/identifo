package grpc

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/grpc/shared"
)

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(settings model.GRPCSettings) (model.UserStorage, error) {
	var err error
	err = nil

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
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("user-storage")
	if err != nil {
		return nil, err
	}

	user := raw.(*shared.GRPCClient)

	user.PluginClient = client

	return user, err
}
