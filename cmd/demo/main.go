package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "./apps.json"
	usersImportPath = "./users.json"
)

var port string

func initServer() model.Server {
	srv, err := boltdb.NewServer(server.ServerSettings)
	if err != nil {
		log.Fatal(err)
	}

	port = server.ServerSettings.GetPort()

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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := initServer()
	log.Println("Demo Identifo server started")
	log.Fatal(http.ListenAndServe(port, s.Router()))
}
