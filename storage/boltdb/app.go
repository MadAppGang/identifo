package boltdb

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/rs/xid"
	bolt "go.etcd.io/bbolt"
)

const (
	// AppBucket is a name for bucket with apps.
	AppBucket = "Apps"
)

// NewAppStorage creates new BoltDB AppStorage implementation.
func NewAppStorage(
	logger *slog.Logger,
	settings model.BoltDBDatabaseSettings) (model.AppStorage, error) {
	if len(settings.Path) == 0 {
		return nil, fmt.Errorf("unable to find init boltdb storage with empty database path")
	}

	// init database
	db, err := InitDB(settings.Path)
	if err != nil {
		return nil, err
	}

	as := AppStorage{
		logger: logger,
		db:     db,
	}
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
	logger *slog.Logger
	db     *bolt.DB
}

// AppByID returns app from memory by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	res := model.AppData{}
	if err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		app := ab.Get([]byte(id))
		if app == nil {
			return model.ErrorNotFound
		}

		var err error
		res, err = model.AppDataFromJSON(app)
		return err
	}); err != nil {
		return model.AppData{}, err
	}
	return res, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return model.AppData{}, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return model.AppData{}, err
	}

	if !app.Active {
		return model.AppData{}, ErrorInactiveApp
	}

	return app, nil
}

// CreateApp creates new app in BoltDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	if len(app.ID) == 0 {
		app.ID = xid.New().String()
	}
	return app, as.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(app)
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))
		return ab.Put([]byte(app.ID), data)
	})
}

// DisableApp disables app in the storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	app.Active = false
	_, err := as.CreateApp(app)
	return err
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	// use ID from the request if it's not set
	if len(newApp.ID) == 0 {
		newApp.ID = appID
	}

	err := as.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(newApp)
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))
		if err := ab.Delete([]byte(appID)); err != nil {
			return err
		}

		return ab.Put([]byte(newApp.ID), data)
	})
	if err != nil {
		return model.AppData{}, err
	}

	updatedApp, err := as.AppByID(newApp.ID)
	return updatedApp, err
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string) ([]model.AppData, error) {
	apps := []model.AppData{}

	err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))

		if iterErr := ab.ForEach(func(k, v []byte) error {
			if strings.Contains(strings.ToLower(string(k)), strings.ToLower(filterString)) {
				app, err := model.AppDataFromJSON(v)
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
		return []model.AppData{}, err
	}
	return apps, nil
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
			_, err := model.AppDataFromJSON(v)
			return err
		})
	})
	return err
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte, cleanOldData bool) error {
	if cleanOldData {
		if err := as.db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte(AppBucket))
			if _, err := tx.CreateBucketIfNotExists([]byte(AppBucket)); err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	apd := []model.AppData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		as.logger.Error("unmarshalling app data",
			logging.FieldError, err)
		return err
	}
	for _, a := range apd {
		if _, err := as.CreateApp(a); err != nil {
			return err
		}
	}
	return nil
}

// Close closes underlying database.
func (as *AppStorage) Close() {
	if err := CloseDB(as.db); err != nil {
		as.logger.Error("Error closing app storage", logging.FieldError, err)
	}
}
