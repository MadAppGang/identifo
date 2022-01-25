package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/grpc/shared"
)

func main() {
	path := flag.String("path", "", "path to database")
	flag.Parse()

	s, err := boltdb.NewUserStorage(model.BoltDBDatabaseSettings{
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

	for {
		<-osch
		s.Close()
		log.Println("Boltdb user storage is terminated.")
		return
	}
}
