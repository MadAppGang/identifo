package mem

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/madappgang/identifo/model"
)

// NewAppStorage creates new in-memory AppStorage implementation.
func NewAppStorage() model.AppStorage {
	return &AppStorage{storage: make(map[string]AppData)}
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	storage map[string]AppData
}

// AppByID returns app by ID from the in-memory storage.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	a, ok := as.storage[id]
	if !ok {
		return nil, ErrorNotFound
	}
	return a, nil
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

//AddNewApp add new app to memory storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	as.storage[app.ID()] = NewAppData(app)
	return app, nil
}

// DisableApp deletes app from in-memory storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID())
	return nil
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	delete(as.storage, oldAppID)
	as.storage[newApp.ID()] = NewAppData(newApp)
	return nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, error) {
	apps := []model.AppData{}
	for _, app := range as.storage {
		if strings.Contains(strings.ToLower(app.Name()), strings.ToLower(filterString)) {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.AddNewApp(AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

type appData struct {
	ID                   string   `json:"id,omitempty"`
	Secret               string   `json:"secret,omitempty"`
	Active               bool     `json:"active,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Offline              bool     `json:"offline,omitempty"`
	RedirectURL          string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan int64    `json:"refresh_token_lifespan,omitempty"`
	TokenLifespan        int64    `json:"token_lifespan,omitempty"`
	TokenPayload         []string `json:"token_payload,omitempty"`
}

// AppData is an in-memory model for model.AppData.
type AppData struct {
	appData
}

// NewAppData instantiates app data in-memory model from the general one.
func NewAppData(data model.AppData) AppData {
	return AppData{appData: appData{
		ID:                   data.ID(),
		Secret:               data.Secret(),
		Active:               data.Active(),
		Name:                 data.Name(),
		Description:          data.Description(),
		Scopes:               data.Scopes(),
		Offline:              data.Offline(),
		RedirectURL:          data.RedirectURL(),
		RefreshTokenLifespan: data.RefreshTokenLifespan(),
		TokenLifespan:        data.TokenLifespan(),
		TokenPayload:         data.TokenPayload(),
	}}
}

// MakeAppData creates new in-memory app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64, tokenPayload []string) AppData {
	return AppData{appData: appData{
		ID:                   id,
		Secret:               secret,
		Active:               active,
		Name:                 name,
		Description:          description,
		Scopes:               scopes,
		Offline:              offline,
		RedirectURL:          redirectURL,
		RefreshTokenLifespan: refreshTokenLifespan,
		TokenLifespan:        tokenLifespan,
		TokenPayload:         tokenPayload,
	}}
}

// ID implements model.AppData interface.
func (ad AppData) ID() string { return ad.appData.ID }

// Secret implements model.AppData interface.
func (ad AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad AppData) Offline() bool { return ad.appData.Offline }

// RedirectURL implements model.AppData interface.
func (ad AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// AddAppDataFromFile loads app data from JSON file and saves it to the storage.
func AddAppDataFromFile(as model.AppStorage, file string) {

}
