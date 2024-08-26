package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/plugin/shared"
)

type wproxy struct {
}

func (w wproxy) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		wproxy{},
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))

	path := flag.String("path", "", "path to database")
	flag.Parse()

	s, err := boltdb.NewUserStorage(
		logging.DefaultLogger,
		model.BoltDBDatabaseSettings{
			Path: *path,
		})

	if err != nil {
		panic(err)
	}

	defer s.Close()

	go plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"user-storage": &shared.UserStoragePlugin{Impl: s},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})

	osch := make(chan os.Signal, 1)
	signal.Notify(osch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-osch
	s.Close()
	logging.DefaultLogger.Info("Boltdb user storage is terminated.")
}
