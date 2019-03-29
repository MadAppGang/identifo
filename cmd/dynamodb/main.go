package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/dynamodb"
)

const (
	testAppID       = "59fd884d8f6b180001f5b4e2"
	appsImportPath  = "../import/apps.json"
	usersImportPath = "../import/users.json"
)

var port string

func initServer() model.Server {
	settings := dynamodb.ServerSettings

	server, err := dynamodb.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	port = settings.GetPort()

	if _, err = server.AppStorage().AppByID(testAppID); err != nil {
		log.Println("Error getting app storage:", err)
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
	s := initServer()
	log.Println("DynamoDB server started")
	log.Fatal(http.ListenAndServe(port, s.Router()))
}
