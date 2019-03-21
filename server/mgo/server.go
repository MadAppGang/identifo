package mgo

import (
	"log"
	"os"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// Settings are the extended settings for MongoDB server.
type Settings struct {
	model.ServerSettings
	DBEndpoint string
	DBName     string
}

const databaseConfigPath = "./database_config.yaml"

// ServerSettings are default server settings.
var ServerSettings = Settings{}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get database configuration file:", err)
	}

	server.LoadConfiguration(dir, databaseConfigPath, &ServerSettings)
	ServerSettings.ServerSettings = server.ServerSettings
}

// NewServer creates new backend service with MongoDB support.
func NewServer(settings Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings.ServerSettings, dbComposer, options...)
}
