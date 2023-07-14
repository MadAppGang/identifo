package mem

import (
	"encoding/json"
	"strings"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/rs/xid"
)

var appNotFoundError = l.NewError(l.ErrorNotFound, "app")

// NewAppStorage creates new in-memory AppStorage implementation.
func NewAppStorage() (model.AppStorage, error) {
	return &AppStorage{storage: make(map[string]model.AppData)}, nil
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	storage map[string]model.AppData
}

// AppByID returns app by ID from the in-memory storage.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	a, ok := as.storage[id]
	if !ok {
		return model.AppData{}, appNotFoundError
	}
	return a, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return model.AppData{}, appNotFoundError
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return model.AppData{}, err
	}

	if !app.Active {
		return model.AppData{}, l.ErrorAPPInactive
	}

	return app, nil
}

// CreateApp creates new app in memory.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	// generate new ID if it's not set
	if len(app.ID) == 0 {
		app.ID = xid.New().String()
	}
	as.storage[app.ID] = app
	return app, nil
}

// DisableApp deletes app from in-memory storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID)
	return nil
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	delete(as.storage, appID)
	as.storage[appID] = newApp
	return newApp, nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string) ([]model.AppData, error) {
	apps := []model.AppData{}
	var total int

	for _, app := range as.storage {
		total++
		if strings.Contains(strings.ToLower(app.Name), strings.ToLower(filterString)) {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

// DeleteApp does nothing here.
func (as *AppStorage) DeleteApp(id string) error {
	return nil
}

// TestDatabaseConnection is always optimistic about the database connection.
func (as *AppStorage) TestDatabaseConnection() error {
	return nil
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte, cleanOldData bool) error {
	if cleanOldData {
		as.storage = make(map[string]model.AppData)
	}

	apd := []model.AppData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		return err
	}
	for _, a := range apd {
		if _, err := as.CreateApp(a); err != nil {
			return err
		}
	}
	return nil
}

// Close clears storage.
func (as *AppStorage) Close() {
	for k := range as.storage {
		delete(as.storage, k)
	}
}
