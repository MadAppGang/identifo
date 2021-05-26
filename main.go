package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	configStoreEtcd "github.com/madappgang/identifo/configuration/storage/etcd"
	configStoreFile "github.com/madappgang/identifo/configuration/storage/file"
	configWatcherEtcd "github.com/madappgang/identifo/configuration/watcher/etcd"
	configWatcherGeneric "github.com/madappgang/identifo/configuration/watcher/generic"
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
	etcdKeyName := flag.String("etcd_key", "identifo", "Key for config settings in etcd folder")
	flag.Parse()

	forever := make(chan struct{}, 1)

	configStorage, err := server.InitConfigurationStorage(*configFlag, *etcdKeyName)
	if err != nil {
		log.Printf("Unable to init config using\n\tconfig string: %s\n\tetcdKeyName: %s\n\twith error: %v\n",
			*configFlag,
			*etcdKeyName,
			err,
		)
		// Trying to fall back to default settings:
		log.Printf("Trying to load default settings from env variable 'SERVER_CONFIG_PATH' or default pathes.\n")
		configStorage, err = configStoreFile.NewDefaultConfigurationStorage()
		if err != nil {
			log.Fatalf("Unable to load default config with error: %v", err)
		}
	}

	srv := initServer(configStorage)
	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
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
	var settings model.ServerSettings
	if err := configStorage.LoadServerSettings(&settings); err != nil {
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

func initWatcher(httpSrv *http.Server, srv model.Server) model.ConfigurationWatcher {
	var cw model.ConfigurationWatcher
	var err error

	watchChan := make(chan interface{}, 1)
	configStorage := srv.ConfigurationStorage()

	switch srv.Settings().ConfigurationStorage.Type {
	case model.ConfigurationStorageTypeEtcd:
		etcdStorage, ok := configStorage.(*configStoreEtcd.ConfigurationStorage)
		if !ok {
			log.Panicln("Incorrect configuration storage type")
		}
		cw, err = configWatcherEtcd.NewConfigurationWatcher(etcdStorage, srv.Settings().ConfigurationStorage.SettingsKey, watchChan)
	case model.ConfigurationStorageTypeS3, model.ConfigurationStorageTypeFile:
		cw, err = configWatcherGeneric.NewConfigurationWatcher(configStorage, srv.Settings().ConfigurationStorage.SettingsKey, watchChan)
	default:
		log.Panicln("Unknown config storage type:", srv.Settings().ConfigurationStorage.Type)
	}

	if err != nil {
		log.Panicln("Cannot init configuration watcher: ", err)
	}

	cw.Watch()
	log.Printf("Watcher initialized (type %s)\n", srv.Settings().ConfigurationStorage.Type)

	go func() {
		for event := range cw.WatchChan() {
			log.Printf("New event from watcher: %+v\n", event)
			var settings model.ServerSettings
			if err := configStorage.LoadServerSettings(&settings); err != nil {
				log.Panicln("Cannot reload server configuration: ", err)
			}

			if err := httpSrv.Shutdown(context.Background()); err != nil {
				log.Panicln("Cannot shutdown server: ", err)
			}

			srv.Close()
			srv = initServer(configStorage)
			*httpSrv = http.Server{Addr: srv.Settings().GetPort()}

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
