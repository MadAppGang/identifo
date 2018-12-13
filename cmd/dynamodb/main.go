package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/dynamodb"
)

func initDB() model.Server {
	settings := dynamodb.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"

	server, err := dynamodb.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		server.ImportApps("../import/apps.json")
		server.ImportUsers("../import/users.json")
	}
	return server
}

func main() {
	r := initDB()
	fmt.Println("DynamoDB server started")
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}
