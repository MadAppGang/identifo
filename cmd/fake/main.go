package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/server/fake"
)

// This server works only with in-memory storages and generated data.
// It should be used for test and CI environments only.
func main() {
	settings := fake.DefaultSettings

	server, err := fake.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(settings.GetPort(), server.Router()))
}
