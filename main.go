package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/madappgang/identifo/config"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
	"github.com/madappgang/identifo/server/dynamodb"
	"github.com/madappgang/identifo/server/fake"
	"github.com/madappgang/identifo/server/mgo"
	"github.com/rs/cors"
)

func main() {
	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
	flag.Parse()

	// ignore error to fall back to default if needed
	settings, _ := model.ConfigStorageSettingsFromString(*configFlag)
	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil {
		log.Printf("Unable to init config using\n\tconfig string: %s\n\twith error: %v\n",
			*configFlag,
			err,
		)
		// Trying to fall back to default settings:
		log.Printf("Trying to load default settings from env variable 'SERVER_CONFIG_PATH' or default pathes.\n")
		configStorage, err = config.DefaultStorage()
		if err != nil {
			log.Fatalf("Unable to load default config with error: %v", err)
		}
	}

	srv := initServer(configStorage)
	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}

	go startHTTPServer(httpSrv)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	log.Println("shutting down the service ⤵️")

	// Stop the service gracefully.
	// TODO: implement gracefull server shutdown
	// srv.Shutdown()
	httpSrv.Shutdown(context.Background())
	log.Println("the server is gracefully stopped, bye 👋")
}

func startHTTPServer(httpSrv *http.Server) {
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func initServer(configStorage model.ConfigurationStorage) model.Server {
	settings, err := configStorage.LoadServerSettings(true)
	if err != nil {
		log.Panicln("Cannot load server settings: ", err)
	}

	dbTypes := make(map[model.DatabaseType]bool)
	var partialComposers []server.PartialDatabaseComposer

	dbTypes[settings.Storage.AppStorage.Type] = true
	dbTypes[settings.Storage.UserStorage.Type] = true
	dbTypes[settings.Storage.TokenStorage.Type] = true
	dbTypes[settings.Storage.TokenBlacklist.Type] = true
	dbTypes[settings.Storage.VerificationCodeStorage.Type] = true
	dbTypes[settings.Storage.InviteStorage.Type] = true

	for dbType := range dbTypes {
		pc, err := initPartialComposer(dbType, settings.Storage)
		if err != nil {
			log.Panicf("Cannot init partial composer for db type %s: %s\n", dbType, err)
		}
		partialComposers = append(partialComposers, pc)
	}

	dbComposer, err := server.NewComposer(settings, partialComposers)
	if err != nil {
		log.Panicln("Cannot init database composer:", err)
	}

	srv, err := server.NewServer(settings, dbComposer, configStorage, &model.CorsOptions{
		API: &cors.Options{AllowedHeaders: []string{"*", "x-identifo-clientid"}, AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}},
	})
	if err != nil {
		log.Panicln("Cannot init server:", err)
	}

	return srv
}

func initPartialComposer(dbType model.DatabaseType, settings model.StorageSettings) (server.PartialDatabaseComposer, error) {
	switch dbType {
	case model.DBTypeBoltDB:
		return boltdb.NewPartialComposer(settings)
	case model.DBTypeMongoDB:
		return mgo.NewPartialComposer(settings)
	case model.DBTypeDynamoDB:
		return dynamodb.NewPartialComposer(settings)
	case model.DBTypeFake:
		return fake.NewPartialComposer(settings)
	}
	return nil, fmt.Errorf("Unknown db type: %s", dbType)
}
