package mem

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

// NewAppStorage creates new in-memory AppStorage implementation.
func NewAppStorage() (model.AppStorage, error) {
	return &AppStorage{storage: make(map[string]AppData)}, nil
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	storage map[string]AppData
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app by ID from the in-memory storage.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	a, ok := as.storage[id]
	if !ok {
		return nil, ErrorNotFound
	}
	return &a, nil
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

// CreateApp creates new app in memory.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(res)
	return result, err
}

// addNewApp adds new app to in-memory storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(*AppData)
	if !ok || a == nil {
		return nil, model.ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if len(a.ID()) == 0 {
		a.appData.ID = xid.New().String()
	}
	as.storage[a.ID()] = NewAppData(app)
	return a, nil
}

// DisableApp deletes app from in-memory storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID())
	return nil
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	delete(as.storage, appID)
	updatedApp := NewAppData(newApp)
	as.storage[newApp.ID()] = updatedApp
	return &updatedApp, nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	apps := []model.AppData{}
	var total int

	for _, app := range as.storage {
		total++
		skip--
		if skip > -1 {
			continue
		}
		if limit != 0 && len(apps) == limit {
			break
		}
		if strings.Contains(strings.ToLower(app.Name()), strings.ToLower(filterString)) {
			apps = append(apps, &app)
		}
	}
	return apps, total, nil
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
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(&AppData{appData: a}); err != nil {
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
