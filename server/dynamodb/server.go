package dynamodb

import (
	"log"
	"os"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

const databaseConfigPath = "./database-config.yaml"

// ServerSettings are server settings.
var ServerSettings = server.ServerSettings

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get database configuration file:", err)
	}

	server.LoadConfiguration(dir, databaseConfigPath, &ServerSettings.DBSettings)
}

// NewServer creates new backend service with DynamoDB support.
func NewServer(settings model.ServerSettings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer, options...)
}
