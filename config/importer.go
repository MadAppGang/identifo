package config

import (
	"io"
	"os"

	"github.com/madappgang/identifo/v2/model"
)

// Importer has a set of helper functions to import data
// Now it uses to import dummy data for isolated tests
// But could be used to import static data for stateless deployment
// more details in docs

// ImportApps imports apps from file.
func ImportApps(filename string, storage model.AppStorage, cleanOldData bool) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return storage.ImportJSON(data, cleanOldData)
}

// ImportUsers imports users from file.
func ImportUsers(filename string, storage model.UserStorage, cleanOldData bool) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return storage.ImportJSON(data, cleanOldData)
}

// ImportUsers imports users from file.
func ImportManagement(filename string, storage model.ManagementKeysStorage, cleanOldData bool) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return storage.ImportJSON(data, cleanOldData)
}

func dataFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
