package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/fake"
)

// This server works only with in-memory storages and generated data.
// It should be used for test and CI environments only.
func main() {
	srv, err := fake.NewServer(server.ServerSettings, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(server.ServerSettings.GetPort(), srv.Router()))
}
