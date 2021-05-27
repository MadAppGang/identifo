package main

import (
	"flag"
	"log"
	"net/http"

	configStoreFile "github.com/madappgang/identifo/configuration/storage/file"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/fake"
)

func loadConfig() model.ConfigurationStorage {
	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
	flag.Parse()

	configStorage, err := server.InitConfigurationStorage(*configFlag)
	if err != nil {
		log.Printf("Unable to init config using\n config string: %s\nwith error: %v\n",
			*configFlag,
			err,
		)
		// Trying to fall back to default settings:
		log.Printf("Trying to load default settings from env variable 'SERVER_CONFIG_PATH' or default pathes.\n")
		configStorage, err = configStoreFile.NewDefaultConfigurationStorage()
		if err != nil {
			log.Fatalf("Unable to load default config with error: %v", err)
		}
	}
	return configStorage
}

func loadServerSettings() model.ServerSettings {
	cs := loadConfig()
	var settings model.ServerSettings
	if err := cs.LoadServerSettings(&settings); err != nil {
		log.Panicln("Cannot load server settings: ", err)
	}
	return settings
}

// This server works only with in-memory storages and generated data.
// It should be used for test and CI environments only.
func main() {
	srv, err := fake.NewServer(loadServerSettings(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(srv.Settings().GetPort(), srv.Router()))
}
