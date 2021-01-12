package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/mongo"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "../import/apps.json"
	usersImportPath = "../import/users.json"
)

func initServer(plugins shared.Plugins) model.Server {
	srv, err := mongo.NewServer(server.ServerSettings, plugins, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = srv.AppStorage().AppByID(testAppID); err != nil {
		log.Println("Error getting app storage:", err)
		if err = srv.ImportApps(appsImportPath); err != nil {
			log.Println("Error importing apps:", err)
		}
		if err = srv.ImportUsers(usersImportPath); err != nil {
			log.Println("Error importing users:", err)
		}
	}
	return srv
}

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

	s := initServer(plugins)
	defer s.Close()

	log.Println("MongoDB server started")
	log.Fatal(http.ListenAndServe(server.ServerSettings.GetPort(), s.Router()))
}
