package main

import (
	"fmt"
	"log"
	"net/http"

	etcdStorage "github.com/madappgang/identifo/configuration/etcd/storage"
	etcdWatcher "github.com/madappgang/identifo/configuration/etcd/watcher"
	mockWatcher "github.com/madappgang/identifo/configuration/mock/watcher"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
	"github.com/madappgang/identifo/server/dynamodb"
	"github.com/madappgang/identifo/server/fake"
	"github.com/madappgang/identifo/server/mgo"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "cmd/import/apps.json"
	usersImportPath = "cmd/import/users.json"
)

func main() {
	forever := make(chan struct{}, 1)

	srv := initServer(nil)
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
	dbTypes := make(map[model.DatabaseType]bool)
	var partialComposers []server.PartialDatabaseComposer

	dbTypes[server.ServerSettings.Storage.AppStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.UserStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.VerificationCodeStorage.Type] = true

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

	configStorageOption := server.ConfigurationStorageOption(configStorage)

	srv, err := server.NewServer(server.ServerSettings, dbComposer, configStorageOption)
	if err != nil {
		log.Panicln("Cannot init server:", err)
	}

	if _, err = srv.AppStorage().AppByID(testAppID); err != nil {
		log.Println("Error getting app by ID:", err)
		if err = srv.ImportApps(appsImportPath); err != nil {
			log.Println("Error importing apps:", err)
		}
		if err = srv.ImportUsers(usersImportPath); err != nil {
			log.Println("Error importing users:", err)
		}
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
		etcdStorage, ok := configStorage.(*etcdStorage.ConfigurationStorage)
		if !ok {
			log.Panicln("Incorrect configuration storage type")
		}
		cw, err = etcdWatcher.NewConfigurationWatcher(etcdStorage, server.ServerSettings.ConfigurationStorage.SettingsKey, watchChan)
	case model.ConfigurationStorageTypeMock:
		cw, err = mockWatcher.NewConfigurationWatcher()
	default:
		log.Panicln("Unknown config storage type:", server.ServerSettings.ConfigurationStorage)
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

			if err := httpSrv.Close(); err != nil {
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
