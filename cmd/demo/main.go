package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/embedded"
)

func initDB() model.Server {
	settings := embedded.DefaultSettings
	settings.StaticFolderPath = "./static"
	settings.PEMFolderPath = "./pem"
	settings.Issuer = "http://localhost:8080"
	settings.MailService = model.MailServiceAWS
	settings.EmailTemplatesPath = "./email_templates"

	server, err := embedded.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		server.ImportApps("./apps.json")
		server.ImportUsers("./users.json")
	}
	return server
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Demo Identifo server started")
	r := initDB()
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}
