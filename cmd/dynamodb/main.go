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
	settings.StaticFolderPath = "../../static"
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
	settings.MailService = model.MailServiceAWS
	settings.EmailTemplatesPath = "../../email_templates"

	server, err := dynamodb.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2"); err != nil {
		log.Println("Error getting app storage:", err)
		if err = server.ImportApps("../import/apps.json"); err != nil {
			log.Println("Error importing apps:", err)
		}
		if err = server.ImportUsers("../import/users.json"); err != nil {
			log.Println("Error importing users:", err)
		}
	}
	return server
}

func main() {
	r := initDB()
	fmt.Println("DynamoDB server started")
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}
