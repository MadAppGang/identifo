package config

import (
	"io/ioutil"
	"os"

	"github.com/madappgang/identifo/v2/model"
)

// Importer has a set of helper functions to import data
// Now it uses to import dummy data for isolated tests
// But could be used to import static data for stateless deployment
// more details in docs

// ImportApps imports apps from file.
func ImportApps(filename string, storage model.AppStorage) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return storage.ImportJSON(data)
}

// ImportUsers imports users from file.
func ImportUsers(filename string, storage model.UserStorage) error {
	data, err := dataFromFile(filename)
	if err != nil {
		return err
	}
	return storage.ImportJSON(data)
}

func dataFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
