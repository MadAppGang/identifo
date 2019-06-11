package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/embedded"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "./apps.json"
	usersImportPath = "./users.json"
)

var port string

func initServer() model.Server {
	settings := embedded.ServerSettings

	server, err := embedded.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	port = settings.GetPort()

	if _, err = server.AppStorage().AppByID(testAppID); err != nil {
		log.Println("Error getting app by ID:", err)
		if err = server.ImportApps(appsImportPath); err != nil {
			log.Println("Error importing apps:", err)
		}
		if err = server.ImportUsers(usersImportPath); err != nil {
			log.Println("Error importing users:", err)
		}
	}

	return server
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := initServer()
	log.Println("Demo Identifo server started")
	log.Fatal(http.ListenAndServe(port, s.Router()))
}
