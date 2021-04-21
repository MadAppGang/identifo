package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	configStoreEtcd "github.com/madappgang/identifo/configuration/storage/etcd"
	configWatcherEtcd "github.com/madappgang/identifo/configuration/watcher/etcd"
	configWatcherGeneric "github.com/madappgang/identifo/configuration/watcher/generic"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
	"github.com/madappgang/identifo/server/dynamodb"
	"github.com/madappgang/identifo/server/fake"
	"github.com/madappgang/identifo/server/mgo"
)

func main() {
	forever := make(chan struct{}, 1)

	configStorage, err := server.InitConfigurationStorage(server.ServerSettings.ConfigurationStorage, server.ServerSettings.StaticFilesStorage.ServerConfigPath)
	if err != nil {
		log.Fatal("Cannot init config storage:", err)
	}

	srv := initServer(configStorage)
	httpSrv := &http.Server{
		Addr:    server.ServerSettings.GetPort(),
		Handler: srv.Router(),
	}

	watcher := initWatcher(httpSrv, srv)
	defer watcher.Stop()

	go startHTTPServer(httpSrv)

	<-forever
}

func startHTTPServer(httpSrv *http.Server) {
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func initServer(configStorage model.ConfigurationStorage) model.Server {
	if err := configStorage.LoadServerSettings(&server.ServerSettings); err != nil {
		log.Panicln("Cannot load server settings: ", err)
	}

	dbTypes := make(map[model.DatabaseType]bool)
	var partialComposers []server.PartialDatabaseComposer

	dbTypes[server.ServerSettings.Storage.AppStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.UserStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenBlacklist.Type] = true
	dbTypes[server.ServerSettings.Storage.VerificationCodeStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.InviteStorage.Type] = true

	for dbType := range dbTypes {
		pc, err := initPartialComposer(dbType, server.ServerSettings.Storage)
		if err != nil {
			log.Panicf("Cannot init partial composer for db type %s: %s\n", dbType, err)
		}
		partialComposers = append(partialComposers, pc)
	}

	dbComposer, err := server.NewComposer(server.ServerSettings, partialComposers)
	if err != nil {
		log.Panicln("Cannot init database composer:", err)
	}

	srv, err := server.NewServer(server.ServerSettings, dbComposer, configStorage, nil)
	if err != nil {
		log.Panicln("Cannot init server:", err)
	}

	return srv
}

func initWatcher(httpSrv *http.Server, srv model.Server) model.ConfigurationWatcher {
	var cw model.ConfigurationWatcher
	var err error

	watchChan := make(chan interface{}, 1)
	configStorage := srv.ConfigurationStorage()

	switch server.ServerSettings.ConfigurationStorage.Type {
	case model.ConfigurationStorageTypeEtcd:
		etcdStorage, ok := configStorage.(*configStoreEtcd.ConfigurationStorage)
		if !ok {
			log.Panicln("Incorrect configuration storage type")
		}
		cw, err = configWatcherEtcd.NewConfigurationWatcher(etcdStorage, server.ServerSettings.ConfigurationStorage.SettingsKey, watchChan)
	case model.ConfigurationStorageTypeS3, model.ConfigurationStorageTypeFile:
		cw, err = configWatcherGeneric.NewConfigurationWatcher(configStorage, server.ServerSettings.ConfigurationStorage.SettingsKey, watchChan)
	default:
		log.Panicln("Unknown config storage type:", server.ServerSettings.ConfigurationStorage.Type)
	}

	if err != nil {
		log.Panicln("Cannot init configuration watcher: ", err)
	}

	cw.Watch()
	log.Printf("Watcher initialized (type %s)\n", server.ServerSettings.ConfigurationStorage.Type)

	go func() {
		for event := range cw.WatchChan() {
			log.Printf("New event from watcher: %+v\n", event)
			if err := configStorage.LoadServerSettings(&server.ServerSettings); err != nil {
				log.Panicln("Cannot reload server configuration: ", err)
			}

			if err := httpSrv.Shutdown(context.Background()); err != nil {
				log.Panicln("Cannot shutdown server: ", err)
			}

			*httpSrv = http.Server{Addr: server.ServerSettings.GetPort()}

			srv.Close()
			srv = initServer(configStorage)

			httpSrv.Handler = srv.Router()

			log.Println("Starting new web server...")
			go startHTTPServer(httpSrv)
		}
	}()
	return cw
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
