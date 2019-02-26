package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/fake"
)

// This server works only with in-memory storages and generated data.
// It should be used for test and CI environments only.
func main() {
	settings := fake.DefaultSettings
	settings.StaticFolderPath = "../../static"
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
	settings.MailService = model.MailServiceAWS
	settings.EmailTemplatesPath = "../../email_templates"

	server, err := fake.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		panic(err)
	}
}
