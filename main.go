package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-plugin"
	configStoreEtcd "github.com/madappgang/identifo/configuration/storage/etcd"
	configWatcherEtcd "github.com/madappgang/identifo/configuration/watcher/etcd"
	configWatcherGeneric "github.com/madappgang/identifo/configuration/watcher/generic"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
	"github.com/madappgang/identifo/server/dynamodb"
	"github.com/madappgang/identifo/server/fake"
	"github.com/madappgang/identifo/server/mongo"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "cmd/import/apps.json"
	usersImportPath = "cmd/import/users.json"
)

func main() {
	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command("sh", "-c", server.ServerSettings.Storage.UserStorage.Path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via gRPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("user_storage")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	plugins := shared.Plugins{
		UserStorage: raw.(shared.UserStorage),
	}

	configStorage, err := server.InitConfigurationStorage(server.ServerSettings.ConfigurationStorage, server.ServerSettings.StaticFilesStorage.ServerConfigPath)
	if err != nil {
		log.Fatal("Cannot init config storage:", err)
	}

	srv := initServer(configStorage, plugins)
	httpSrv := &http.Server{
		Addr:    server.ServerSettings.GetPort(),
		Handler: srv.Router(),
	}

	watcher := initWatcher(httpSrv, srv, plugins)
	defer watcher.Stop()

	go startHTTPServer(httpSrv)

	defer srv.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Println("Got SIGINT...")
	case syscall.SIGTERM:
		log.Println("Got SIGTERM...")
	}
}

func startHTTPServer(httpSrv *http.Server) {
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func initServer(configStorage model.ConfigurationStorage, plugins shared.Plugins) model.Server {
	if err := configStorage.LoadServerSettings(&server.ServerSettings); err != nil {
		log.Panicln("Cannot load server settings: ", err)
	}

	dbTypes := make(map[model.DatabaseType]bool)
	var partialComposers []server.PartialDatabaseComposer

	dbTypes[server.ServerSettings.Storage.AppStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenBlacklist.Type] = true
	dbTypes[server.ServerSettings.Storage.VerificationCodeStorage.Type] = true

	for dbType := range dbTypes {
		pc, err := initPartialComposer(dbType, server.ServerSettings.Storage, plugins)
		if err != nil {
			log.Panicf("Cannot init partial composer for db type %s: %s\n", dbType, err)
		}
		partialComposers = append(partialComposers, pc)
	}

	dbComposer, err := server.NewComposer(server.ServerSettings, partialComposers, plugins)
	if err != nil {
		log.Panicln("Cannot init database composer:", err)
	}

	srv, err := server.NewServer(server.ServerSettings, dbComposer, configStorage, nil)
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

func initWatcher(httpSrv *http.Server, srv model.Server, plugins shared.Plugins) model.ConfigurationWatcher {
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
			srv = initServer(configStorage, plugins)

			httpSrv.Handler = srv.Router()

			log.Println("Starting new web server...")
			go startHTTPServer(httpSrv)
		}
	}()
	return cw
}

func initPartialComposer(dbType model.DatabaseType, settings model.StorageSettings, plugins shared.Plugins) (server.PartialDatabaseComposer, error) {
	switch dbType {
	case model.DBTypeBoltDB:
		return boltdb.NewPartialComposer(settings, plugins)
	case model.DBTypeMongoDB:
		return mongo.NewPartialComposer(settings, plugins)
	case model.DBTypeDynamoDB:
		return dynamodb.NewPartialComposer(settings, plugins)
	case model.DBTypeFake:
		return fake.NewPartialComposer(settings)
	}
	return nil, fmt.Errorf("Unknown db type: %s", dbType)
}
