package main

import (
	"fmt"
	"log"

	"github.com/akrylysov/algnhsa"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/server/boltdb"
	"github.com/madappgang/identifo/server/dynamodb"
	"github.com/madappgang/identifo/server/fake"
	"github.com/madappgang/identifo/server/mgo"
)

const (
	testAppID       = "testAppID"
	appsImportPath  = "./import/apps.json"
	usersImportPath = "./import/users.json"
)

func main() {
	srv := initServer()
	algnhsa.ListenAndServe(srv.Router(), nil)
}

func initServer() model.Server {
	dbTypes := make(map[model.DatabaseType]bool)
	var partialComposers []server.PartialDatabaseComposer

	dbTypes[server.ServerSettings.Storage.AppStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.UserStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.TokenStorage.Type] = true
	dbTypes[server.ServerSettings.Storage.VerificationCodeStorage.Type] = true

	for dbType := range dbTypes {
		pc, err := initPartialComposer(dbType, server.ServerSettings.Storage)
		if err != nil {
			log.Panicf("Cannot init partial composer for db type %s: %s\n", dbType, err)
		}
		partialComposers = append(partialComposers, pc)
	}

	dbComposer, err := server.NewComposer(server.ServerSettings, partialComposers)
	if err != nil {
		log.Panicln("Cannot init database composer:", err)
	}

	srv, err := server.NewServer(server.ServerSettings, dbComposer)
	if err != nil {
		log.Panicln("Cannot init server:", err)
	}

	if _, err = srv.AppStorage().AppByID(testAppID); err != nil {
		log.Println("Error getting app storage:", err)
		if err = srv.ImportApps(appsImportPath); err != nil {
			log.Println("Error importing apps:", err)
		}
		if err = srv.ImportUsers(usersImportPath); err != nil {
			log.Println("Error importing users:", err)
		}
	}

	return srv
}

func initPartialComposer(dbType model.DatabaseType, settings model.StorageSettings) (server.PartialDatabaseComposer, error) {
	switch dbType {
	case model.DBTypeBoltDB:
		return boltdb.NewPartialComposer(settings)
	case model.DBTypeMongoDB:
		return mgo.NewPartialComposer(settings)
	case model.DBTypeDynamoDB:
		return dynamodb.NewPartialComposer(settings)
	case model.DBTypeFake:
		return fake.NewPartialComposer(settings)
	}
	return nil, fmt.Errorf("Unknown db type: %s", dbType)
}
