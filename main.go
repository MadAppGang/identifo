package main

import (
	"context"
	"log"
	"net/http"

	etcdStorage "github.com/madappgang/identifo/configuration/storage/etcd"
	etcdWatcher "github.com/madappgang/identifo/configuration/watcher/etcd"
	watcherMock "github.com/madappgang/identifo/configuration/watcher/mock"
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

	srv := initServer()
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
	if err := httpSrv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Restarting server...")
		} else {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}
}

func initServer() model.Server {
	var err error
	var dbComposer server.DatabaseComposer

	switch server.ServerSettings.DBType {
	case model.DBTypeBoltDB:
		dbComposer, err = boltdb.NewComposer(server.ServerSettings)
	case model.DBTypeDynamoDB:
		dbComposer, err = dynamodb.NewComposer(server.ServerSettings)
	case model.DBTypeMongoDB:
		dbComposer, err = mgo.NewComposer(server.ServerSettings)
	case model.DBTypeFake:
		dbComposer, err = fake.NewComposer(server.ServerSettings)
	default:
		log.Panicln("Unknown database type:", server.ServerSettings.DBType)
	}
	if err != nil {
		log.Panicln("Cannot init database composer:", err)
	}

	srv, err := server.NewServer(server.ServerSettings, dbComposer)
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

	switch server.ServerSettings.ConfigurationStorage {
	case model.ConfigurationStorageTypeEtcd:
		etcdStorage, ok := configStorage.(*etcdStorage.ConfigurationStorage)
		if !ok {
			log.Panicln("Incorrect configuration storage type")
		}
		cw, err = etcdWatcher.NewConfigurationWatcher(etcdStorage, watchChan)
	case model.ConfigurationStorageTypeMock:
		cw, err = watcherMock.NewConfigurationWatcher(watchChan)
	default:
		log.Panicln("Unknown config storage type:", server.ServerSettings.ConfigurationStorage)
	}

	if err != nil {
		log.Panicln("Cannot init configuration watcher: ", err)
	}

	cw.Watch()
	log.Printf("Watcher initialized (type %s)\n", server.ServerSettings.ConfigurationStorage)

	go func() {
		for event := range cw.WatchChan() {
			log.Printf("New event from watcher: %v+\n", event)
			if err := configStorage.LoadServerSettings(&server.ServerSettings); err != nil {
				log.Panicln("Cannot reload server configuration: ", err)
			}

			if err := httpSrv.Shutdown(context.Background()); err != nil {
				log.Panicln("Cannot shutdown server: ", err)
			}
			go startHTTPServer(&http.Server{
				Addr:    server.ServerSettings.GetPort(),
				Handler: initServer().Router(),
			})
		}
	}()
	return cw
}
