package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/mongo"
	"github.com/madappgang/identifo/v2/storage/plugin/shared"
)

func main() {
	connectionString := flag.String("connection", "", "mongo connection string")
	databaseName := flag.String("database", "", "name of database")
	flag.Parse()

	s, err := mongo.NewUserStorage(model.MongoDatabaseSettings{
		ConnectionString: *connectionString,
		DatabaseName:     *databaseName,
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
	log.Println("Mongo user storage is terminated.")
}
