package plugin

import (
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	grpcShared "github.com/madappgang/identifo/v2/user_payload_provider/grpc/shared"
	"github.com/madappgang/identifo/v2/user_payload_provider/plugin/shared"
)

// NewTokenPayloadProvider creates and inits plugin for payload provider.
func NewTokenPayloadProvider(
	logger *slog.Logger,
	settings model.PluginSettings,
	timeout time.Duration,
) (model.TokenPayloadProvider, error) {
	params := []string{}
	for k, v := range settings.Params {
		params = append(params, "-"+k)
		params = append(params, v)
	}

	cfg := &plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command(settings.Cmd, params...),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logging.NewHCLogger(logger, false),
	}

	if settings.RedirectStd {
		cfg.SyncStdout = os.Stdout
		cfg.SyncStderr = os.Stderr
	}

	client := plugin.NewClient(cfg)

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

	tpp := raw.(*grpcShared.GRPCClient)

	if timeout == 0 {
		timeout = time.Second
	}

	tpp.Timeout = timeout
	tpp.Closable = pluginClosableClient{client: client}

	return tpp, err
}

type pluginClosableClient struct {
	client *plugin.Client
}

func (g pluginClosableClient) Close() error {
	g.client.Kill()
	return nil
}
