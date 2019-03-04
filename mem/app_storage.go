package mem

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/madappgang/identifo/model"
)

//NewAppStorage creates new memory AppStorage implementation
func NewAppStorage() model.AppStorage {
	as := AppStorage{}
	as.storage = make(map[string]AppData)
	return &as
}

//AppStorage is fully functional app storage
type AppStorage struct {
	storage map[string]AppData
}

//AppByID returns app from memory by ID
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	if a, ok := as.storage[id]; !ok {
		return nil, ErrorNotFound
	} else {
		return a, nil
	}
}

//AddNewApp add new app to memory storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	as.storage[app.ID()] = NewAppData(app)
	return app, nil
}

//DisableApp disables app from storage
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID())
	return nil
}

//UpdateApp updates app in storage
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

// DeleteApp does nothing here.
func (as *AppStorage) DeleteApp(id string) error {
	return nil
}

//ImportJSON import data from JSON
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println(err)
		return err
	}
	for _, a := range apd {
		_, err := as.AddNewApp(AppData{appData: a})
		if err != nil {
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

//AppData is memory model for model.AppData
type AppData struct {
	appData
}

//NewAppData instantiate app data memory model from general one
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

//MakeAppData creates new memory app data instance
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

func (ad AppData) ID() string     { return ad.appData.ID }
func (ad AppData) Secret() string { return ad.appData.Secret }
func (ad AppData) Active() bool   { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad AppData) Name() string                { return ad.appData.Name }
func (ad AppData) Description() string         { return ad.appData.Description }
func (ad AppData) Scopes() []string            { return ad.appData.Scopes }
func (ad AppData) Offline() bool               { return ad.appData.Offline }
func (ad AppData) RedirectURL() string         { return ad.appData.RedirectURL }
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }
func (ad AppData) TokenLifespan() int64        { return ad.appData.TokenLifespan }
func (ad AppData) TokenPayload() []string      { return ad.appData.TokenPayload }

//AddAppDataFromFile loads appdata from JSON file and save it to the storage
func AddAppDataFromFile(as model.AppStorage, file string) {

}
