package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/mgo"
)

func server() model.Server {
	settings := mgo.DefaultSettings
	settings.StaticFolderPath = "../../static"
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
	settings.DBEndpoint = "localhost:27017"
	settings.DBName = "identifo"
	settings.MailService = model.MailServiceAWS
	settings.EmailTemplatesPath = "../../email_templates"

	server, err := mgo.NewServer(settings)
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
	fmt.Println("mongodb server started")
	r := server()

	if err := http.ListenAndServe(":8080", r.Router()); err != nil {
		panic(err)
	}
}
