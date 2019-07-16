package boltdb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

const (
	// AppBucket is a name for bucket with apps.
	AppBucket = "Apps"
)

// NewAppStorage creates new BoltDB AppStorage implementation.
func NewAppStorage(db *bolt.DB) (model.AppStorage, error) {
	as := AppStorage{db: db}
	// ensure we have app's bucket in the database
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(AppBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &as, nil
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	db *bolt.DB
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app from memory by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	res := new(AppData)
	if err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		app := ab.Get([]byte(id))
		if app == nil {
			return model.ErrorNotFound
		}

		var err error
		res, err = AppDataFromJSON(app)
		return err
	}); err != nil {
		return nil, err
	}
	return res, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return nil, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return nil, err
	}

	if !app.Active() {
		return nil, ErrorInactiveApp
	}

	return app, nil
}

// CreateApp creates new app in BoltDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(res)
	return result, err
}

// addNewApp adds new app to memory storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if len(res.ID()) == 0 {
		res.appData.ID = xid.New().String()
	}
	return res, as.db.Update(func(tx *bolt.Tx) error {
		data, err := res.Marshal()
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))

		return ab.Put([]byte(res.ID()), data)
	})
}

// DisableApp disables app in the storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return ErrorWrongDataFormat
	}
	res.appData.Active = false
	_, err := as.addNewApp(res)
	return err
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	res, ok := newApp.(*AppData)
	if !ok || newApp == nil {
		return nil, model.ErrorWrongDataFormat
	}
	// use ID from the request if it's not set
	if len(res.ID()) == 0 {
		res.appData.ID = appID
	}

	err := as.db.Update(func(tx *bolt.Tx) error {
		data, err := res.Marshal()
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))
		if err := ab.Delete([]byte(appID)); err != nil {
			return err
		}

		return ab.Put([]byte(res.ID()), data)
	})
	if err != nil {
		return nil, err
	}

	updatedApp, err := as.AppByID(res.ID())
	return updatedApp, err

}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	apps := []model.AppData{}
	var total int

	err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))

		if iterErr := ab.ForEach(func(k, v []byte) error {
			if strings.Contains(strings.ToLower(string(k)), strings.ToLower(filterString)) {
				total++
				skip--
				if skip > -1 || (limit != 0 && len(apps) == limit) {
					return nil
				}
				app, err := AppDataFromJSON(v)
				if err != nil {
					return err
				}
				apps = append(apps, app)
			}
			return nil
		}); iterErr != nil {
			return iterErr
		}
		return nil
	})

	if err != nil {
		return []model.AppData{}, 0, err
	}
	return apps, total, nil
}

// DeleteApp deletes app by ID.
func (as *AppStorage) DeleteApp(id string) error {
	err := as.db.Update(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		return ab.Delete([]byte(id))
	})
	return err
}

// TestDatabaseConnection checks whether we can fetch the first document in the applications bucket.
func (as *AppStorage) TestDatabaseConnection() error {
	err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		return ab.ForEach(func(k, v []byte) error {
			_, err := AppDataFromJSON(v)
			return err
		})
	})
	return err
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println(err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(&AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

// Close closes underlying database.
func (as *AppStorage) Close() {
	if err := as.db.Close(); err != nil {
		log.Printf("Error closing app storage: %s\n", err)
	}
}
